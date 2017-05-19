package bolt

/*

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/theopticians/optician-api/core"
	"github.com/theopticians/optician-api/core/store"
)

func TestBoltStore(t *testing.T) {
	core.GenericTestStore(t, func() store.Store {
		return NewBoltStore("optician_test_" + RandStringBytes(10) + ".db")
	})

	removeBoltDatabases(".")
}

func removeBoltDatabases(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		match, _ := regexp.MatchString(".*.db", name)
		if match {
			err = os.Remove(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

*/
