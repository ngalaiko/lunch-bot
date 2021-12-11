package lunch

import (
	"testing"
	"time"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

func TestHistory_roll_boost__active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())
	placeNames := []string{"place1", "place2", "place3"}
	places := make([]*places.Place, len(placeNames))
	for i, name := range placeNames {
		place, err := roller.NewPlace(ctx, name)
		assertNoError(t, err)
		places[i] = place
	}

	_, _, firstRollError := roller.Roll(ctx, today)
	assertNoError(t, firstRollError)

	_, firstBoostError := roller.Boost(ctx, places[0].ID, today.Add(time.Minute))
	assertNoError(t, firstBoostError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 1, len(history.ActiveBoosts[places[0].ID]))

	_, anotherBoostError := roller.Boost(testContext(testUser()), places[0].ID, today.Add(time.Minute))
	assertNoError(t, anotherBoostError)

	history, err = roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 2, len(history.ActiveBoosts[places[0].ID]))
}

func TestHistory_boost_roll__no_active_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())
	placeNames := []string{"place1", "place2", "place3"}
	places := make([]*places.Place, len(placeNames))
	for i, name := range placeNames {
		place, err := roller.NewPlace(ctx, name)
		assertNoError(t, err)
		places[i] = place
	}

	_, firstBoostError := roller.Boost(ctx, places[0].ID, today)
	assertNoError(t, firstBoostError)

	_, anotherBoostError := roller.Boost(testContext(testUser()), places[0].ID, today)
	assertNoError(t, anotherBoostError)

	_, _, firstRollError := roller.Roll(ctx, today.Add(time.Minute))
	assertNoError(t, firstRollError)

	history, err := roller.buildHistory(ctx, today.Add(2*time.Minute))
	assertNoError(t, err)
	assertEqual(t, 0, len(history.ActiveBoosts[places[0].ID]))
}
