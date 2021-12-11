package lunch

import (
	"context"
	"fmt"
	"math"
	"time"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/users"
)

const (
	hoursInADay     = 24
	totalPoints     = 1
	boostMultiplier = 5
)

type rollsHistory struct {
	ThisWeekBoosts []*boosts.Boost
	ThisWeekRolls  map[time.Weekday][]*rolls.Roll
	LastRolled     map[places.ID]time.Time
	ActiveBoosts   map[places.ID][]*boosts.Boost
}

func (r *Roller) buildHistory(ctx context.Context, now time.Time) (*rollsHistory, error) {
	allRolls, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	year, week := now.ISOWeek()
	thisWeekRolls := map[time.Weekday][]*rolls.Roll{}
	lastRolled := map[places.ID]time.Time{}
	var latestRoll *rolls.Roll
	for _, roll := range allRolls {
		rollYear, rollWeek := roll.Time.ISOWeek()
		sameYear := rollYear == year
		sameWeek := rollWeek == week
		if sameYear && sameWeek {
			weekday := roll.Time.Weekday()
			thisWeekRolls[weekday] = append(thisWeekRolls[weekday], roll)
		}

		if roll.Time.After(lastRolled[roll.PlaceID]) {
			lastRolled[roll.PlaceID] = roll.Time
		}

		if latestRoll == nil {
			latestRoll = roll
		} else if roll.Time.After(latestRoll.Time) {
			latestRoll = roll
		}
	}

	thisWeekBoosts := []*boosts.Boost{}
	activeBoosts := map[places.ID][]*boosts.Boost{}
	for _, boost := range allBoosts {
		boostYear, boostWeek := boost.Time.ISOWeek()
		sameYear := boostYear == year
		sameWeek := boostWeek == week
		if sameYear && sameWeek {
			thisWeekBoosts = append(thisWeekBoosts, boost)
		}

		// boosts lasts until the next roll
		if latestRoll == nil {
			activeBoosts[boost.PlaceID] = append(activeBoosts[boost.PlaceID], boost)
		} else if latestRoll.Time.Before(boost.Time) {
			activeBoosts[boost.PlaceID] = append(activeBoosts[boost.PlaceID], boost)
		}
	}

	return &rollsHistory{
		ThisWeekRolls:  thisWeekRolls,
		ThisWeekBoosts: thisWeekBoosts,
		LastRolled:     lastRolled,
		ActiveBoosts:   activeBoosts,
	}, nil
}

func (h *rollsHistory) CanBoost(user *users.User, now time.Time) error {
	if h.pointsLeft(user) <= 0 {
		return ErrNoPoints
	}

	return nil
}

func (h *rollsHistory) CanRoll(user *users.User, now time.Time) error {
	firstRollToday := len(h.ThisWeekRolls[now.Weekday()]) == 0
	if firstRollToday {
		// anyone can make the first roll a day
		return nil
	}

	if h.pointsLeft(user) <= 0 {
		return ErrNoPoints
	}

	return nil
}

func (h *rollsHistory) pointsLeft(user *users.User) int {
	points := totalPoints

	for _, boost := range h.ThisWeekBoosts {
		if boost.UserID == user.ID {
			// Boost costs one point
			points--
		}
	}

	for _, rolls := range h.ThisWeekRolls {
		if len(rolls) <= 1 {
			// first roll a day is always allowed
			continue
		}
		for _, roll := range rolls[1:] {
			// consecutive rolls a day are rerolls, only one reroll per week is allowed
			if roll.UserID == user.ID {
				points--
			}
		}
	}
	return points
}

// getWeights returns a list of weights for places to choose from.
// higher weights means higher chance of choosing a place.
// list of weights is the same length as list of places,
// where weights[i] is the weight for places[i].
//
// weight are distributed in a way so that the most recent rolls get the lowest weight.
func (h *rollsHistory) getWeights(places []*places.Place, now time.Time) []float64 {
	placesTotal := len(places)
	weights := make([]float64, placesTotal)
	for i, place := range places {
		lastRolledAt, wasRolled := h.LastRolled[place.ID]
		if !wasRolled {
			weights[i] = float64(placesTotal)
		} else {
			rolledAgo := now.Sub(lastRolledAt)
			rolledDaysAgo := int(math.Floor(rolledAgo.Hours() / hoursInADay))
			if rolledDaysAgo >= placesTotal {
				weights[i] = float64(placesTotal)
			} else {
				weights[i] = float64(rolledDaysAgo) + 1
			}
		}

		for range h.ActiveBoosts[place.ID] {
			weights[i] *= boostMultiplier
		}
	}
	return weights
}