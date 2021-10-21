package lunch

import (
	"testing"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

func TestHistory_roll_boost__active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(store.NewInMemory())
	pp := []string{"place1", "place2", "place3"}
	for _, place := range pp {
		assertNoError(t, roller.NewPlace(ctx, place))
	}

	_, firstRollError := roller.Roll(ctx, today)
	assertNoError(t, firstRollError)

	firstBoostError := roller.Boost(ctx, pp[0], today.Add(time.Minute))
	assertNoError(t, firstBoostError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 1, len(history.ActiveBoosts[places.Name(pp[0])]))

	anotherBoostError := roller.Boost(testContext(testUser()), pp[0], today.Add(time.Minute))
	assertNoError(t, anotherBoostError)

	history, err = roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 2, len(history.ActiveBoosts[places.Name(pp[0])]))
}

func TestHistory_boost_roll__no_active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(store.NewInMemory())
	pp := []string{"place1", "place2", "place3"}
	for _, place := range pp {
		assertNoError(t, roller.NewPlace(ctx, place))
	}

	firstBoostError := roller.Boost(ctx, pp[0], today)
	assertNoError(t, firstBoostError)

	anotherBoostError := roller.Boost(testContext(testUser()), pp[0], today)
	assertNoError(t, anotherBoostError)

	_, firstRollError := roller.Roll(ctx, today.Add(time.Minute))
	assertNoError(t, firstRollError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 0, len(history.ActiveBoosts[places.Name(pp[0])]))
}
