package lunch

import (
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
	ThisWeekBoosts  []*boosts.Boost
	RollsPerWeekday map[time.Weekday][]*rolls.Roll
	LastRolled      map[places.ID]time.Time
	ActiveBoosts    map[places.ID]int
}

func buildHistory(allRolls map[rolls.ID]*rolls.Roll, allBoosts map[boosts.ID]*boosts.Boost, now time.Time) *rollsHistory {
	year, week := now.ISOWeek()
	rollsPerWeekday := map[time.Weekday][]*rolls.Roll{}
	lastRolled := map[places.ID]time.Time{}
	var latestRoll *rolls.Roll
	for _, roll := range allRolls {
		rollYear, rollWeek := roll.Time.ISOWeek()
		sameYear := rollYear == year
		sameWeek := rollWeek == week
		if sameYear && sameWeek {
			weekday := roll.Time.Weekday()
			rollsPerWeekday[weekday] = append(rollsPerWeekday[weekday], roll)
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
	activeBoosts := map[places.ID]int{}
	for _, boost := range allBoosts {
		boostYear, boostWeek := boost.Time.ISOWeek()
		sameYear := boostYear == year
		sameWeek := boostWeek == week
		if sameYear && sameWeek {
			thisWeekBoosts = append(thisWeekBoosts, boost)
		}

		// boosts lasts until the next roll
		if latestRoll == nil {
			activeBoosts[boost.PlaceID]++
		} else if latestRoll.Time.Before(boost.Time) {
			activeBoosts[boost.PlaceID]++
		}
	}

	return &rollsHistory{
		RollsPerWeekday: rollsPerWeekday,
		ThisWeekBoosts:  thisWeekBoosts,
		LastRolled:      lastRolled,
		ActiveBoosts:    activeBoosts,
	}
}

func (h *rollsHistory) CanBoost(userID users.ID, now time.Time) error {
	if h.pointsLeft(userID) <= 0 {
		return ErrNoPoints
	}

	return nil
}

func (h *rollsHistory) CanRoll(userID users.ID, now time.Time) error {
	firstRollToday := len(h.RollsPerWeekday[now.Weekday()]) == 0
	if firstRollToday {
		// anyone can make the first roll a day
		return nil
	}

	if h.pointsLeft(userID) <= 0 {
		return ErrNoPoints
	}

	return nil
}

func (h *rollsHistory) pointsLeft(userID users.ID) int {
	points := totalPoints

	for _, boost := range h.ThisWeekBoosts {
		if boost.UserID == userID {
			// Boost costs one point
			points--
		}
	}

	for _, rolls := range h.RollsPerWeekday {
		if len(rolls) <= 1 {
			// first roll a day is always allowed
			continue
		}

		for _, roll := range rolls[1:] {
			// consecutive rolls a day are rerolls, only one reroll per week is allowed
			if roll.UserID == userID {
				points--
			}
		}
	}

	return points
}

// getWeights returns a list of weights for places to choose from.
// higher weights means higher chance of choosing a place.
//
// weight are distributed in a way so that the most recent rolls get the lowest weight.
func (h *rollsHistory) getWeights(allPlaces map[places.ID]*places.Place, now time.Time) map[places.ID]float64 {
	placesTotal := len(allPlaces)
	weights := make(map[places.ID]float64, placesTotal)
	for placeID, place := range allPlaces {
		lastRolledAt, wasRolled := h.LastRolled[place.ID]
		if !wasRolled {
			weights[placeID] = float64(placesTotal)
		} else {
			rolledAgo := now.Sub(lastRolledAt)
			rolledDaysAgo := int(math.Floor(rolledAgo.Hours() / hoursInADay))
			if rolledDaysAgo >= placesTotal {
				weights[placeID] = float64(placesTotal)
			} else {
				weights[placeID] = float64(rolledDaysAgo) + 1
			}
		}

		for i := 0; i < h.ActiveBoosts[place.ID]; i++ {
			weights[placeID] *= boostMultiplier
		}
	}
	return weights
}
