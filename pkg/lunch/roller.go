package lunch

import (
	"context"
	"fmt"
	"math"
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

const (
	hoursInADay        = 24
	defaultProbability = 1.0
)

type Roller struct {
	placesStore *places.Store
	rollsStore  *rolls.Store

	rand *rand.Rand
}

func New(storage store.Storage) *Roller {
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

	if err := r.placesStore.Store(ctx, places.NewPlace(name, user)); err != nil {
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

func (r *Roller) Roll(ctx context.Context, now time.Time) (*places.Place, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}

	rollsHistory, err := r.rollsStore.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	if err := r.checkRules(ctx, user, rollsHistory, now); err != nil {
		return nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	place, err := r.pickRandomPlace(ctx, rollsHistory, now)
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

func (r *Roller) pickRandomPlace(ctx context.Context, rollsHistory []*rolls.Roll, now time.Time) (*places.Place, error) {
	allNames, err := r.placesStore.ListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	weights := getWeights(allNames, rollsHistory, now)
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

// getWeights returns a list of weights for places to choose from.
// higher weights means higher chance of choosing a place.
// list of weights is the same length as list of allNames,
// where weights[i] is the weight for allNames[i].
//
// weight are distributed in a way so that the most recent rolls get the lowest weight.
func getWeights(allNames []string, rollsHistory []*rolls.Roll, now time.Time) []float64 {
	lastRolled := map[string]time.Time{}
	for _, roll := range rollsHistory {
		if roll.Time.After(lastRolled[roll.PlaceName]) {
			lastRolled[roll.PlaceName] = roll.Time
		}
	}

	namesTotal := len(allNames)
	weights := make([]float64, namesTotal)
	for i, name := range allNames {
		lastRolledAt, wasRolled := lastRolled[name]
		if !wasRolled {
			weights[i] = defaultProbability
		} else {
			rolledAgo := now.Sub(lastRolledAt)
			rolledDaysAgo := int(math.Floor(rolledAgo.Hours() / hoursInADay))
			if rolledDaysAgo >= namesTotal {
				weights[i] = defaultProbability
			} else {
				weights[i] = defaultProbability / float64(namesTotal-rolledDaysAgo)
			}
		}
	}

	return weights
}

func (r *Roller) checkRules(ctx context.Context, user *users.User, rollsHistory []*rolls.Roll, now time.Time) error {
	year, week := now.ISOWeek()
	rollsByWeekday := map[time.Weekday][]*rolls.Roll{}
	for _, roll := range rollsHistory {
		rollYear, rollWeek := roll.Time.ISOWeek()
		sameYear := rollYear == year
		sameWeek := rollWeek == week
		if sameYear && sameWeek {
			weekday := roll.Time.Weekday()
			rollsByWeekday[weekday] = append(rollsByWeekday[weekday], roll)
		}
	}

	firstRollToday := len(rollsByWeekday[now.Weekday()]) == 0
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
