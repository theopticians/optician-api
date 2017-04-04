package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"

	"github.com/boltdb/bolt"
)

var errNotFound = errors.New("Key not found in DB")

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string) Store {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte("results"))
		_, err = tx.CreateBucketIfNotExists([]byte("images"))
		_, err = tx.CreateBucketIfNotExists([]byte("baseimages"))
		_, err = tx.CreateBucketIfNotExists([]byte("masks"))
		return err
	})

	if err != nil {
		panic(err)
	}

	return &BoltStore{db}
}

func (s *BoltStore) Close() {
	s.db.Close()
}

func (s *BoltStore) generateUniqueKey(projectID, branch, target, browser string) string {
	return projectID + "|" + branch + "|" + target + "|" + browser
}

func (s *BoltStore) getStringValue(bucket, key string) (string, error) {
	b, err := s.getValue(bucket, key)

	if err != nil {
		return "", err
	}

	return string(b), err
}

func (s *BoltStore) getValue(bucket, key string) ([]byte, error) {
	var results []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		results = b.Get([]byte(key))
		return nil
	})

	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errNotFound
	}

	return results, nil
}

func (s *BoltStore) storeValue(bucket, key string, value []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return b.Put([]byte(key), value)
	})
	return err
}

func (s *BoltStore) GetTestList() []string {
	ret := []string{}
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("results"))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ret = append(ret, string(k))
		}

		return nil
	})

	return ret
}

func (s *BoltStore) StoreResults(r Results) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("results"))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(r)
		if err != nil {
			return err
		}
		return b.Put([]byte(r.TestID), encoded)
	})
	return err

}

func (s *BoltStore) GetResults(ID string) (Results, error) {
	val, err := s.getValue("results", ID)

	res := Results{}

	if err != nil {
		return res, err
	}

	err = json.Unmarshal(val, &res)

	return res, err
}

func (s *BoltStore) GetMasks(projectID, branch, target, browser string) ([]image.Rectangle, error) {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	serialized, err := s.getValue("masks", key)
	if err != nil {
		if err == errNotFound {
			return []image.Rectangle{}, nil
		}
		return nil, err
	}

	var ret []image.Rectangle

	err = json.Unmarshal(serialized, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *BoltStore) StoreMasks(masks []image.Rectangle, projectID, branch, target, browser string) error {
	key := s.generateUniqueKey(projectID, branch, target, browser)

	serialized, err := json.Marshal(masks)
	if err != nil {
		return err
	}

	return s.storeValue("masks", key, serialized)
}

func (s *BoltStore) GetImage(imgID string) (image.Image, error) {
	val, err := s.getValue("images", imgID)

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(val)

	img, _, err := image.Decode(buff)

	return img, err
}

func (s *BoltStore) StoreImage(img image.Image) (string, error) {
	imgID := RandStringBytes(10)

	buffer := bytes.NewBuffer(nil)
	err := png.Encode(buffer, img)
	if err != nil {
		return "", err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("images"))
		if err != nil {
			return err
		}

		return b.Put([]byte(imgID), buffer.Bytes())
	})

	return imgID, err
}

func (s *BoltStore) GetBaseImageID(projectID, branch, target, browser string) (string, error) {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	return s.getStringValue("baseimages", key)
}

func (s *BoltStore) SetBaseImageID(baseImageID, projectID, branch, target, browser string) error {
	key := s.generateUniqueKey(projectID, branch, target, browser)

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("baseimages"))
		if err != nil {
			return err
		}

		println(key, baseImageID)

		return b.Put([]byte(key), []byte(baseImageID))
	})

	return err
}
