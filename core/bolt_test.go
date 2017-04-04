package core

import "testing"

func TestBoltStore(t *testing.T) {
	testStore(t, func() Store {
		return NewBoltStore("optician_test.db")
	})
}
