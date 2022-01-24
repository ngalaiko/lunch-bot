package lunch

import (
	"context"
	"fmt"
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

type Place struct {
	*places.Place
	User   *users.User `json:"user"`
	Chance float64     `json:"chance"`
}

type Boost struct {
	*boosts.Boost
	User  *users.User   `json:"user"`
	Place *places.Place `json:"place"`
}

type Roll struct {
	*rolls.Roll
	User  *users.User   `json:"user"`
	Place *places.Place `json:"place"`
}

type views struct {
	placesStore  storage_places.Storage
	rollsStore   storage_rolls.Storage
	boostsStore  storage_boosts.Storage
	usersStorage storage_users.Storage
}

func (v *views) Boosts(ctx context.Context, bb map[boosts.ID]*boosts.Boost) ([]*Boost, error) {
	if len(bb) == 0 {
		return nil, nil
	}

	allPlaces, err := v.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	allUsers, err := v.usersStorage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	boosts := make([]*Boost, 0, len(bb))
	for _, b := range bb {
		boosts = append(boosts, &Boost{
			Boost: b,
			User:  allUsers[b.UserID],
			Place: allPlaces[b.PlaceID],
		})
	}

	return boosts, nil
}

func (v *views) Rolls(ctx context.Context, rr map[rolls.ID]*rolls.Roll) ([]*Roll, error) {
	if len(rr) == 0 {
		return nil, nil
	}

	allPlaces, err := v.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	allUsers, err := v.usersStorage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	rolls := make([]*Roll, 0, len(rr))
	for _, r := range rr {
		rolls = append(rolls, &Roll{
			Roll:  r,
			User:  allUsers[r.UserID],
			Place: allPlaces[r.PlaceID],
		})
	}

	return rolls, nil
}

func (v *views) Places(ctx context.Context, now time.Time, pp map[places.ID]*places.Place) ([]*Place, error) {
	if len(pp) == 0 {
		return nil, nil
	}

	allBoosts, err := v.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	allRolls, err := v.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allPlaces, err := v.placesStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	allUsers, err := v.usersStorage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	history, err := buildHistory(allRolls, allBoosts, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build history")
	}

	weights := history.getWeights(pp, now)
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	views := make([]*Place, 0, len(allPlaces))
	for i, weight := range weights {
		if _, ok := pp[i]; !ok {
			continue
		}
		chance := weight / weightsSum
		views = append(views, &Place{
			Place:  pp[i],
			User:   allUsers[pp[i].UserID],
			Chance: chance,
		})
	}

	return views, nil
}
