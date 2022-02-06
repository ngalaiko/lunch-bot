package store

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestListWithValues(t *testing.T) {
	dir := os.TempDir()
	dbFile, err := ioutil.TempFile(dir, "bolt-")
	assertNoError(t, err)

	bolt, err := NewBolt(dbFile.Name())
	assertNoError(t, err)

	type value struct {
		A string
	}
	v := value{"a"}
	assertNoError(t, bolt.Put(context.Background(), "bucket", "key", v))

	var dest []value
	err = bolt.List(context.Background(), "bucket", &dest)
	assertNoError(t, err)

	t.Log(dest)
}

func TestListWithPtr(t *testing.T) {
	dir := os.TempDir()
	dbFile, err := ioutil.TempFile(dir, "bolt-")
	assertNoError(t, err)

	bolt, err := NewBolt(dbFile.Name())
	assertNoError(t, err)

	type value struct {
		A string
	}
	v := &value{"a"}
	assertNoError(t, bolt.Put(context.Background(), "bucket", "key", v))

	var dest []*value
	err = bolt.List(context.Background(), "bucket", &dest)
	assertNoError(t, err)

	t.Log(dest)
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	assertError(t, nil, err)
}

func assertError(t *testing.T, expected error, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}
