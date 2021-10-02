package lunch

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/store"
	"lunch/pkg/users"
)

var (
	ErrNoRerolls = fmt.Errorf("no rerolls left")
	ErrNoPlaces  = fmt.Errorf("no places to choose from")
)

type Roller struct {
	placesStore *places.Store
	rollsStore  *rolls.Store

	rand *rand.Rand
}

func New(storage *store.S3Store) *Roller {
	return &Roller{
		placesStore: places.NewStore(storage),
		rollsStore:  rolls.NewStore(storage),
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) NewPlace(ctx context.Context, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	if err := r.placesStore.Add(ctx, places.NewPlace(name, user)); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}
	return nil
}

func (r *Roller) ListPlacesNames(ctx context.Context) ([]string, error) {
	names, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}
	return names, nil
}

func (r *Roller) Roll(ctx context.Context) (*places.Place, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	if err := r.checkRules(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	place, err := r.pickRandomPlace(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to pick random place: %w", err)
	}

	if err := r.storeResult(ctx, user, place); err != nil {
		return nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	return place, nil
}

func (r *Roller) storeResult(ctx context.Context, user *users.User, place *places.Place) error {
	return r.rollsStore.Store(ctx, rolls.NewRoll(user, place.Name))
}

func (r *Roller) pickRandomPlace(ctx context.Context) (*places.Place, error) {
	names, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places ids: %w", err)
	}

	if len(names) == 0 {
		return nil, ErrNoPlaces
	}

	randomPlaceID := names[r.rand.Intn(len(names))]
	place, err := r.placesStore.Get(ctx, randomPlaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get place by id: %w", err)
	}

	return place, nil
}

func (r *Roller) checkRules(ctx context.Context, user *users.User) error {
	rollsThisWeek, err := r.rollsStore.ListThisWeekRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list this week rolls: %w", err)
	}

	rollsByWeekday := map[time.Weekday][]*rolls.Roll{}
	for _, roll := range rollsThisWeek {
		weekday := roll.Time.Weekday()
		rollsByWeekday[weekday] = append(rollsByWeekday[weekday], roll)
	}

	today := time.Now()
	firstRollToday := len(rollsByWeekday[today.Weekday()]) == 0
	if firstRollToday {
		// anyone can make the first roll a day
		return nil
	}

	for _, rolls := range rollsByWeekday {
		if len(rolls) <= 1 {
			continue
		}
		for _, roll := range rolls[1:] {
			if roll.UserID == user.ID {
				// consecutive rolls a day are rerolls, only one reroll per week is allowed
				return ErrNoRerolls
			}
		}
	}

	return nil
}
