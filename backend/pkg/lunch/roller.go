package lunch

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	"lunch/pkg/lunch/rolls"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/users"
	storage_users "lunch/pkg/users/storage"
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
	views  *views

	rand *rand.Rand
}

func New(
	placesStorage storage_places.Storage,
	boostsStorage storage_boosts.Storage,
	rollsStorage storage_rolls.Storage,
	eventsRegistry *events.Registry,
	userService storage_users.Storage,
) *Roller {
	return &Roller{
		placesStore: placesStorage,
		rollsStore:  rollsStorage,
		boostsStore: boostsStorage,
		events:      eventsRegistry,
		views: &views{
			placesStore:  placesStorage,
			boostsStore:  boostsStorage,
			rollsStore:   rollsStorage,
			usersStorage: userService,
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) OnRollCreated(fn func(context.Context, *Roll) error) {
	r.events.OnRollCreated(func(ctx context.Context, roll *rolls.Roll) error {
		rr, err := r.views.Rolls(ctx, map[rolls.ID]*rolls.Roll{roll.ID: roll})
		if err != nil {
			return fmt.Errorf("failed to view roll: %w", err)
		}
		return fn(ctx, rr[0])
	})
}

func (r *Roller) OnPlaceCreated(fn func(context.Context, *Place) error) {
	r.events.OnPlaceCreated(func(ctx context.Context, place *places.Place) error {
		pp, err := r.views.Places(ctx, time.Now(), map[places.ID]*places.Place{place.ID: place})
		if err != nil {
			return fmt.Errorf("failed to view place: %w", err)
		}
		return fn(ctx, pp[0])
	})
}

func (r *Roller) OnBoostCreated(fn func(context.Context, *Boost) error) {
	r.events.OnBoostCreated(func(ctx context.Context, boost *boosts.Boost) error {
		bb, err := r.views.Boosts(ctx, map[boosts.ID]*boosts.Boost{boost.ID: boost})
		if err != nil {
			return fmt.Errorf("failed to view boost: %w", err)
		}
		return fn(ctx, bb[0])
	})
}

func (r *Roller) GetPlace(ctx context.Context, id places.ID) (*places.Place, error) {
	return r.placesStore.GetByID(ctx, id)
}

func (r *Roller) NewPlace(ctx context.Context, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	place := places.NewPlace(places.Name(name), user)
	if err := r.placesStore.Store(ctx, place); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}

	r.events.PlaceCreated(place)

	return nil
}

func (r *Roller) ListRolls(ctx context.Context) ([]*Roll, error) {
	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}
	return r.views.Rolls(ctx, allRolls)
}

func (r *Roller) ListBoosts(ctx context.Context) ([]*Boost, error) {
	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}
	return r.views.Boosts(ctx, allBoosts)
}

func (r *Roller) ListPlaces(ctx context.Context, now time.Time) ([]*Place, error) {
	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	return r.views.Places(ctx, now, allPlaces)
}

func (r *Roller) Boost(ctx context.Context, placeID places.ID, now time.Time) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list boosts: %w", err)
	}

	history, err := buildHistory(allRolls, allBoosts, now)
	if err != nil {
		return fmt.Errorf("failed to build history: %w", err)
	}

	if err := history.CanBoost(user, now); err != nil {
		return fmt.Errorf("can't boost any more: %w", err)
	}

	boost := boosts.NewBoost(user, placeID, now)
	if err := r.boostsStore.Store(ctx, boost); err != nil {
		return fmt.Errorf("failed to store boost: %w", err)
	}

	r.events.BoostCreated(boost)

	return nil
}

func (r *Roller) Roll(ctx context.Context, now time.Time) (*Roll, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	history, err := buildHistory(allRolls, allBoosts, now)
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

	roll := rolls.NewRoll(user, place.ID, now)
	if err := r.rollsStore.Store(ctx, roll); err != nil {
		return nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	r.events.RollCreated(roll)

	rr, err := r.views.Rolls(ctx, map[rolls.ID]*rolls.Roll{roll.ID: roll})
	if err != nil {
		return nil, err
	}

	return rr[0], nil
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
	return randomPlace, nil
}

// weightedRandom returns a random index i from the slice of weights, proportional to the weights[i] value.
func weightedRandom(rand *rand.Rand, weights map[places.ID]float64) places.ID {
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
	return ""
}
