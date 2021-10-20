package lunch

import (
	"testing"
	"time"

	"lunch/pkg/store"
)

func TestHistory_roll_boost__active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(store.NewInMemory())
	places := []string{"place1", "place2", "place3"}
	for _, place := range places {
		assertNoError(t, roller.NewPlace(ctx, place))
	}

	_, firstRollError := roller.Roll(ctx, today)
	assertNoError(t, firstRollError)

	firstBoostError := roller.Boost(ctx, places[0], today.Add(time.Minute))
	assertNoError(t, firstBoostError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)

	assertEqual(t, places[0], string(history.ActiveBoost.PlaceName))
}

func TestHistory_boost_roll__no_active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(store.NewInMemory())
	places := []string{"place1", "place2", "place3"}
	for _, place := range places {
		assertNoError(t, roller.NewPlace(ctx, place))
	}

	firstBoostError := roller.Boost(ctx, places[0], today)
	assertNoError(t, firstBoostError)

	_, firstRollError := roller.Roll(ctx, today.Add(time.Minute))
	assertNoError(t, firstRollError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)

	assertNil(t, history.ActiveBoost)
}
