package storage

import (
	"context"
	"errors"
	"io/ioutil"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/store"
	"lunch/pkg/users"
	"reflect"
	"testing"
)

func Test_PlaceDeleted(t *testing.T) {
	file, err := ioutil.TempFile("", "test-bolt")
	assertNoError(t, err)
	bolt, err := store.NewBolt(file.Name())
	assertNoError(t, err)

	storage := New(events.NewBoltStorage(bolt))

	place := places.NewPlace(rooms.ID("1"), users.ID("1"), "test")
	assertNoError(t, storage.Create(context.Background(), place))

	fromDB, err := storage.Place(context.Background(), place.RoomID, place.ID)
	assertNoError(t, err)
	assertEqual(t, place.ID, fromDB.ID)

	assertNoError(t, storage.Delete(context.Background(), users.ID("2"), place))

	_, err2 := storage.Place(context.Background(), place.RoomID, place.ID)
	assertError(t, err2, ErrNotFound)
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

func assertError(t *testing.T, expected error, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}
