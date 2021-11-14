package lunch

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
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

	rand *rand.Rand
}

func New(
	placesStorage storage_places.Storage,
	boostsStorage storage_boosts.Storage,
	rollsStorage storage_rolls.Storage,
) *Roller {
	return &Roller{
		placesStore: placesStorage,
		rollsStore:  rollsStorage,
		boostsStore: boostsStorage,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) NewPlace(ctx context.Context, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	if err := r.placesStore.Store(ctx, places.NewPlace(places.Name(name), user)); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}

	return nil
}

func (r *Roller) ListPlaces(ctx context.Context) ([]places.Name, error) {
	names, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}
	return names, nil
}

func (r *Roller) ListChances(ctx context.Context, now time.Time) (map[places.Name]float64, error) {
	allNames, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allNames) == 0 {
		return nil, ErrNoPlaces
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build history")
	}

	weights := history.getWeights(allNames, now)
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	chances := make(map[places.Name]float64, len(allNames))
	for i, weight := range weights {
		chance := weight / weightsSum * 100
		chances[allNames[i]] = chance
	}

	return chances, nil
}

func (r *Roller) Boost(ctx context.Context, name string, now time.Time) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to build history: %w", err)
	}

	if err := history.CanBoost(user, now); err != nil {
		return fmt.Errorf("can't boost any more: %w", err)
	}

	boost := boosts.NewBoost(user, places.Name(name), now)
	if err := r.boostsStore.Store(ctx, boost); err != nil {
		return fmt.Errorf("failed to store boost: %w", err)
	}

	return nil
}

func (r *Roller) Roll(ctx context.Context, now time.Time) (*places.Place, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	history, err := r.buildHistory(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build history: %w", err)
	}

	if err := history.CanRoll(user, now); err != nil {
		return nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	place, err := r.pickRandomPlace(ctx, history, now)
	if err != nil {
		return nil, fmt.Errorf("failed to pick random place: %w", err)
	}

	if err := r.storeResult(ctx, user, place, now); err != nil {
		return nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	return place, nil
}

func (r *Roller) storeResult(ctx context.Context, user *users.User, place *places.Place, now time.Time) error {
	return r.rollsStore.Store(ctx, rolls.NewRoll(user, place.Name, now))
}

func (r *Roller) pickRandomPlace(ctx context.Context, history *rollsHistory, now time.Time) (*places.Place, error) {
	allNames, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allNames) == 0 {
		return nil, ErrNoPlaces
	}

	weights := history.getWeights(allNames, now)
	randomIndex := weightedRandom(r.rand, weights)
	randomPlaceName := allNames[randomIndex]
	place, err := r.placesStore.GetByName(ctx, randomPlaceName)
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
