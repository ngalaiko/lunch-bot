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
	storage_users "lunch/pkg/users/storage"
)

var (
	ErrNoPoints = fmt.Errorf("no points left")
	ErrNoPlaces = fmt.Errorf("no places to choose from")
)

type Roller struct {
	*registry

	placesStore storage_places.Storage
	rollsStore  storage_rolls.Storage
	boostsStore storage_boosts.Storage
	usersStore  storage_users.Storage

	rand *rand.Rand
}

func New(
	placesStorage storage_places.Storage,
	boostsStorage storage_boosts.Storage,
	rollsStorage storage_rolls.Storage,
	userStorage storage_users.Storage,
) *Roller {
	return &Roller{
		registry:    newEventsRegistry(),
		placesStore: placesStorage,
		rollsStore:  rollsStorage,
		boostsStore: boostsStorage,
		usersStore:  userStorage,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) CreatePlace(ctx context.Context, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	place := places.NewPlace(user.ID, name)
	if err := r.placesStore.Store(ctx, place); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}

	r.PlaceCreated(&Place{
		Place: place,
		User:  user,
	})

	return nil
}

func (r *Roller) ListRolls(ctx context.Context) ([]*Roll, error) {
	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}
	if len(allRolls) == 0 {
		return nil, nil
	}

	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	rolls := make([]*Roll, 0, len(allRolls))
	for _, roll := range allRolls {
		rolls = append(rolls, &Roll{
			Roll:  roll,
			User:  allUsers[roll.UserID],
			Place: allPlaces[roll.PlaceID],
		})
	}

	return rolls, nil
}

func (r *Roller) ListBoosts(ctx context.Context) ([]*Boost, error) {
	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	boosts := make([]*Boost, 0, len(allBoosts))
	for _, b := range allBoosts {
		boosts = append(boosts, &Boost{
			Boost: b,
			User:  allUsers[b.UserID],
			Place: allPlaces[b.PlaceID],
		})
	}

	return boosts, nil
}

func (r *Roller) ListPlaces(ctx context.Context, now time.Time) ([]*Place, error) {
	allPlaces, err := r.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	history := buildHistory(allRolls, allBoosts, now)
	weights := history.getWeights(allPlaces, now)
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	views := make([]*Place, 0, len(allPlaces))
	for i, weight := range weights {
		if _, ok := allPlaces[i]; !ok {
			continue
		}

		chance := weight / weightsSum
		views = append(views, &Place{
			Place:  allPlaces[i],
			User:   allUsers[allPlaces[i].UserID],
			Chance: chance,
		})
	}

	return views, nil
}

func (r *Roller) CreateBoost(ctx context.Context, placeID places.ID, now time.Time) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	place, err := r.placesStore.GetByID(ctx, placeID)
	if err != nil {
		return fmt.Errorf("failed to get place: %w", err)
	}

	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list boosts: %w", err)
	}

	history := buildHistory(allRolls, allBoosts, now)

	if err := history.CanBoost(user.ID, now); err != nil {
		return fmt.Errorf("can't boost any more: %w", err)
	}

	boost := boosts.NewBoost(user.ID, placeID, now)
	if err := r.boostsStore.Store(ctx, boost); err != nil {
		return fmt.Errorf("failed to store boost: %w", err)
	}

	r.BoostCreated(&Boost{
		Boost: boost,
		User:  user,
		Place: place,
	})

	return nil
}

func (r *Roller) CreateRoll(ctx context.Context, now time.Time) (*Roll, error) {
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

	history := buildHistory(allRolls, allBoosts, now)
	if err := history.CanRoll(user.ID, now); err != nil {
		return nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	place, err := r.pickRandomPlace(ctx, history, now)
	if err != nil {
		return nil, fmt.Errorf("failed to pick random place: %w", err)
	}

	roll := rolls.NewRoll(user.ID, place.ID, now)
	if err := r.rollsStore.Store(ctx, roll); err != nil {
		return nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	rollView := &Roll{
		Roll:  roll,
		User:  user,
		Place: place,
	}

	r.RollCreated(rollView)

	return rollView, nil
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

	panic("should never happen")
}
