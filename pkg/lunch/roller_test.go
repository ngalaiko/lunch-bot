package lunch

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/users"
)

func TestRoll_noPlaces(t *testing.T) {
	t.Parallel()

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())

	place, err := roller.Roll(ctx, time.Now())
	assertError(t, ErrNoPlaces, err)
	assertNil(t, place)
}

func TestRoll_reroll_then_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())
	placeNames := []string{"place1", "place2", "place3"}
	places := make([]*places.Place, len(placeNames))
	for i, name := range placeNames {
		place, err := roller.NewPlace(ctx, name)
		assertNoError(t, err)
		places[i] = place
	}

	_, firstRollError := roller.Roll(ctx, today)
	assertNoError(t, firstRollError)

	_, firstRerollError := roller.Roll(ctx, today.Add(1*time.Minute))
	assertNoError(t, firstRerollError)

	firstBoostError := roller.Boost(ctx, places[0].ID, today.Add(2*time.Minute))
	assertError(t, ErrNoPoints, firstBoostError)

	nextWeekBoostError := roller.Boost(ctx, places[0].ID, today.Add(oneWeek))
	assertNoError(t, nextWeekBoostError)
}

func TestRoll_boost_then_reroll(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())
	placeNames := []string{"place1", "place2", "place3"}
	places := make([]*places.Place, len(placeNames))
	for i, name := range placeNames {
		place, err := roller.NewPlace(ctx, name)
		assertNoError(t, err)
		places[i] = place
	}

	_, firstRollError := roller.Roll(ctx, today)
	assertNoError(t, firstRollError)

	firstBoostError := roller.Boost(ctx, places[0].ID, today.Add(1*time.Minute))
	assertNoError(t, firstBoostError)

	secondBoostError := roller.Boost(ctx, places[0].ID, today.Add(2*time.Minute))
	assertError(t, ErrNoPoints, secondBoostError)

	_, firstRerollError := roller.Roll(ctx, today)
	assertError(t, ErrNoPoints, firstRerollError)

	nextWeekBoostError := roller.Boost(ctx, places[0].ID, today.Add(oneWeek))
	assertNoError(t, nextWeekBoostError)
}

func TestRoll_rerolls(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay

	user1 := testUser()
	rolls := []struct {
		Description string
		By          *users.User
		When        time.Time
		Error       error
	}{
		{"first roll today",
			user1, today, nil},
		{"second roll today - first reroll this week",
			user1, today.Add(time.Minute), nil},
		{"third roll today - second reroll this week",
			user1, today.Add(2 * time.Minute), ErrNoPoints},
		{"first roll tomorrow",
			user1, today.Add(oneDay), nil},
		{"second roll tomorrow - second reroll this week",
			user1, today.Add(oneDay).Add(time.Minute), ErrNoPoints},
		{"first roll next week - allowed",
			user1, today.Add(oneWeek), nil},
		{"second roll next week - first reroll that week",
			user1, today.Add(oneWeek).Add(time.Minute), nil},
		{"third roll next week - second reroll that week",
			user1, today.Add(oneWeek).Add(2 * time.Minute), ErrNoPoints},
	}

	ctx := testContext(testUser())
	roller := New(storage_places.NewMemory(), storage_boosts.NewMemory(), storage_rolls.NewMemory())
	placeNames := []string{"place1", "place2", "place3"}
	for _, name := range placeNames {
		_, err := roller.NewPlace(ctx, name)
		assertNoError(t, err)
	}

	for _, roll := range rolls {
		t.Run(roll.Description, func(t *testing.T) {
			place, err := roller.Roll(testContext(roll.By), roll.When)
			if roll.Error != nil {
				assertError(t, roll.Error, err)
				assertNil(t, place)
			} else {
				assertNoError(t, err)
				assertNotNil(t, place)
			}
		})
	}
}

var userID *int64 = new(int64)

func testUser() *users.User {
	id := atomic.AddInt64(userID, 1)
	return &users.User{
		ID:   fmt.Sprint(id),
		Name: fmt.Sprintf("test user - %d", id),
	}
}

func testContext(u *users.User) context.Context {
	return users.NewContext(context.Background(), u)
}

func assertNil(t *testing.T, v interface{}) {
	t.Helper()

	if !reflect.ValueOf(v).IsNil() {
		t.Errorf("\nexpected: %+v\ngot: %+v", nil, v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	t.Helper()

	if reflect.ValueOf(v).IsNil() {
		t.Errorf("\nexpected: %+v\ngot: %+v", nil, v)
	}
}

func assertError(t *testing.T, expected error, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	assertError(t, nil, err)
}

func assertEqual(t *testing.T, expected, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}
