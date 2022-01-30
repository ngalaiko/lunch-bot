package lunch

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/store"
	"lunch/pkg/users"
	storage_users "lunch/pkg/users/storage"
)

func TestRoll_noPlaces(t *testing.T) {
	t.Parallel()

	ctx := testContext(testUser())
	file, err := ioutil.TempFile("", "test-bolt")
	assertNoError(t, err)
	bolt, err := store.NewBolt(file.Name())
	assertNoError(t, err)
	roller := New(storage_places.NewBolt(bolt), storage_boosts.NewBolt(bolt), storage_rolls.NewBolt(bolt), storage_users.NewBolt(bolt))

	place, err := roller.CreateRoll(ctx, time.Now())
	assertError(t, ErrNoPlaces, err)
	assertNil(t, place)
}

func TestRoll_reroll_then_boost(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay

	ctx := testContext(testUser())
	file, err := ioutil.TempFile("", "test-bolt")
	assertNoError(t, err)
	bolt, err := store.NewBolt(file.Name())
	assertNoError(t, err)
	roller := New(storage_places.NewBolt(bolt), storage_boosts.NewBolt(bolt), storage_rolls.NewBolt(bolt), storage_users.NewBolt(bolt))
	placeNames := []string{"place1", "place2", "place3"}
	for _, name := range placeNames {
		assertNoError(t, roller.CreatePlace(ctx, name))
	}

	places, err := roller.ListPlaces(ctx, today)
	assertNoError(t, err)

	_, firstRollError := roller.CreateRoll(ctx, today)
	assertNoError(t, firstRollError)

	_, firstRerollError := roller.CreateRoll(ctx, today.Add(1*time.Minute))
	assertNoError(t, firstRerollError)

	firstBoostError := roller.CreateBoost(ctx, places[0].ID, today.Add(2*time.Minute))
	assertError(t, ErrNoPoints, firstBoostError)

	nextWeekBoostError := roller.CreateBoost(ctx, places[0].ID, today.Add(oneWeek))
	assertNoError(t, nextWeekBoostError)
}

func TestRoll_boost_then_reroll(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay

	ctx := testContext(testUser())
	file, err := ioutil.TempFile("", "test-bolt")
	assertNoError(t, err)
	bolt, err := store.NewBolt(file.Name())
	assertNoError(t, err)
	roller := New(storage_places.NewBolt(bolt), storage_boosts.NewBolt(bolt), storage_rolls.NewBolt(bolt), storage_users.NewBolt(bolt))
	placeNames := []string{"place1", "place2", "place3"}
	for _, name := range placeNames {
		assertNoError(t, roller.CreatePlace(ctx, name))
	}

	places, err := roller.ListPlaces(ctx, today)
	assertNoError(t, err)

	_, firstRollError := roller.CreateRoll(ctx, today)
	assertNoError(t, firstRollError)

	firstBoostError := roller.CreateBoost(ctx, places[0].ID, today.Add(1*time.Minute))
	assertNoError(t, firstBoostError)

	secondBoostError := roller.CreateBoost(ctx, places[0].ID, today.Add(2*time.Minute))
	assertError(t, ErrNoPoints, secondBoostError)

	_, firstRerollError := roller.CreateRoll(ctx, today)
	assertError(t, ErrNoPoints, firstRerollError)

	nextWeekBoostError := roller.CreateBoost(ctx, places[0].ID, today.Add(oneWeek))
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
	file, err := ioutil.TempFile("", "test-bolt-*")
	assertNoError(t, err)
	bolt, err := store.NewBolt(file.Name())
	assertNoError(t, err)
	roller := New(storage_places.NewBolt(bolt), storage_boosts.NewBolt(bolt), storage_rolls.NewBolt(bolt), storage_users.NewBolt(bolt))

	placeNames := []string{"place1", "place2", "place3"}
	for _, name := range placeNames {
		assertNoError(t, roller.CreatePlace(ctx, name))
	}

	for _, expected := range rolls {
		t.Run(expected.Description, func(t *testing.T) {
			roll, err := roller.CreateRoll(testContext(expected.By), expected.When)
			if expected.Error != nil {
				assertError(t, expected.Error, err)
				assertNil(t, roll)
			} else {
				assertNoError(t, err)
				assertNotNil(t, roll)
			}
		})
	}
}

var userID *int64 = new(int64)

func testUser() *users.User {
	id := atomic.AddInt64(userID, 1)
	return &users.User{
		ID:   users.ID(fmt.Sprint(id)),
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
