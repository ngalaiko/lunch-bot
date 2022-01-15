package lunch

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	"lunch/pkg/lunch/rolls"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/users"
)

var (
	ErrNoPoints = fmt.Errorf("no points left")
	ErrNoPlaces = fmt.Errorf("no places to choose from")
)

type Roller struct {
	placesStore storage_places.Storage
	rollsStore  storage_rolls.Storage
	boostsStore storage_boosts.Storage

	events *events.Registry

	rand *rand.Rand
}

func New(
	placesStorage storage_places.Storage,
	boostsStorage storage_boosts.Storage,
	rollsStorage storage_rolls.Storage,
	eventsRegistry *events.Registry,
) *Roller {
	return &Roller{
		placesStore: placesStorage,
		rollsStore:  rollsStorage,
		boostsStore: boostsStorage,
		events:      eventsRegistry,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) NewPlace(ctx context.Context, name string) (*places.Place, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	place := places.NewPlace(places.Name(name), user)
	if err := r.placesStore.Store(ctx, place); err != nil {
		return nil, fmt.Errorf("failed to store place: %w", err)
	}

	r.events.PlaceCreated(place)

	return place, nil
}

func (r *Roller) ListRolls(ctx context.Context) ([]*rolls.Roll, error) {
	return r.rollsStore.ListRolls(ctx)
}

type Place struct {
	*places.Place
	Chance float64 `json:"chance"`
}

func (r *Roller) ListPlaces(ctx context.Context, now time.Time) ([]*Place, error) {
	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build history")
	}

	weights := history.getWeights(allPlaces, now)
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	chances := make([]*Place, len(allPlaces))
	for i, weight := range weights {
		chance := weight / weightsSum
		chances[i] = &Place{
			Place:  allPlaces[i],
			Chance: chance,
		}
	}

	sort.SliceStable(chances, func(i, j int) bool {
		return chances[i].Name < chances[j].Name
	})

	sort.SliceStable(chances, func(i, j int) bool {
		return chances[i].Chance > chances[j].Chance
	})

	return chances, nil
}

func (r *Roller) Boost(ctx context.Context, placeID places.ID, now time.Time) (*boosts.Boost, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build history: %w", err)
	}

	if err := history.CanBoost(user, now); err != nil {
		return nil, fmt.Errorf("can't boost any more: %w", err)
	}

	boost := boosts.NewBoost(user, placeID, now)
	if err := r.boostsStore.Store(ctx, boost); err != nil {
		return nil, fmt.Errorf("failed to store boost: %w", err)
	}

	r.events.BoostCreated(boost)

	return boost, nil
}

func (r *Roller) Roll(ctx context.Context, now time.Time) (*rolls.Roll, *places.Place, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("expected to find who in the context")
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build history: %w", err)
	}

	if err := history.CanRoll(user, now); err != nil {
		return nil, nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	place, err := r.pickRandomPlace(ctx, history, now)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pick random place: %w", err)
	}

	roll := rolls.NewRoll(user, place.ID, now)
	if err := r.rollsStore.Store(ctx, roll); err != nil {
		return nil, nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	r.events.RollCreated(roll)

	return roll, place, nil
}

func (r *Roller) pickRandomPlace(ctx context.Context, history *rollsHistory, now time.Time) (*places.Place, error) {
	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	weights := history.getWeights(allPlaces, now)
	randomIndex := weightedRandom(r.rand, weights)
	randomPlace := allPlaces[randomIndex]
	place, err := r.placesStore.GetByID(ctx, randomPlace.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get place by id: %w", err)
	}

	return place, nil
}

// weightedRandom returns a random index i from the slice of weights, proportional to the weights[i] value.
func weightedRandom(rand *rand.Rand, weights []float64) int {
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	remainingDistance := rand.Float64() * weightsSum
	for i, weight := range weights {
		remainingDistance -= weight
		if remainingDistance < 0 {
			return i
		}
	}

	// never happens
	return -1
}
