package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"sort"

	"github.com/boltdb/bolt"
)

var (
	imagesBucket     = []byte("images")
	resultsBucket    = []byte("results")
	baseImagesBucket = []byte("baseImages")
	baseMasksBucket  = []byte("baseMasks")
	masksBucket      = []byte("masks")
)

var KeyNotFoundError = errors.New("Key not found in DB")

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
		_, err = tx.CreateBucketIfNotExists(resultsBucket)
		_, err = tx.CreateBucketIfNotExists(imagesBucket)
		_, err = tx.CreateBucketIfNotExists(baseImagesBucket)
		_, err = tx.CreateBucketIfNotExists(masksBucket)
		_, err = tx.CreateBucketIfNotExists(baseMasksBucket)
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

func (s *BoltStore) getStringValue(bucket []byte, key string) (string, error) {
	b, err := s.getValue(bucket, key)

	if err != nil {
		return "", err
	}

	return string(b), err
}

func (s *BoltStore) getValue(bucket []byte, key string) ([]byte, error) {
	var results []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		results = b.Get([]byte(key))
		return nil
	})

	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, KeyNotFoundError
	}

	return results, nil
}

func (s *BoltStore) storeStringValue(bucket []byte, key, value string) error {
	return s.storeValue(bucket, key, []byte(value))
}

func (s *BoltStore) storeValue(bucket []byte, key string, value []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), value)
	})
	return err
}

func (s *BoltStore) GetResults() ([]Result, error) {
	ret := []Result{}
	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(resultsBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			r := Result{}

			err := json.Unmarshal(v, &r)
			if err != nil {
				return err
			}

			ret = append(ret, r)
		}

		return nil
	})

	sort.Slice(ret, func(i, j int) bool { return ret[i].Timestamp.After(ret[j].Timestamp) })

	return ret, err
}

func (s *BoltStore) GetResultsByBatch(batch string) ([]Result, error) {

	ret := []Result{}
	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(resultsBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			r := Result{}

			err := json.Unmarshal(v, &r)
			if err != nil {
				return err
			}

			if r.Batch == batch {
				ret = append(ret, r)
			}
		}

		return nil
	})

	sort.Slice(ret, func(i, j int) bool { return ret[i].Timestamp.After(ret[j].Timestamp) })

	return ret, err
}

func (s *BoltStore) GetBatchs() ([]BatchInfo, error) {
	ret := []BatchInfo{}

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(resultsBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t Result

			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}

			found := -1
			for i := 0; i < len(ret); i++ {
				if ret[i].ID == t.Batch {
					found = i
					break
				}
			}

			if found >= 0 {
				if t.Timestamp.After(ret[found].Timestamp) {
					ret[found].Timestamp = t.Timestamp
				}

				if t.DiffScore > 0 {
					ret[found].Failed++
				}
			} else {
				failed := 0
				if t.DiffScore > 0 {
					failed++
				}
				ret = append(ret, BatchInfo{t.Batch, t.Timestamp, failed})
			}

		}

		return nil
	})

	return ret, err
}

func (s *BoltStore) GetLastResult(projectID, branch, target, browser string) (Result, error) {
	ret := Result{}
	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(resultsBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t Result

			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}

			if t.ProjectID == projectID && t.Branch == branch && t.Target == target && t.Browser == browser {
				if t.Timestamp.After(ret.Timestamp) {
					ret = t
				}
			}
		}

		return nil
	})

	return ret, err
}

func (s *BoltStore) StoreResult(r Result) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(resultsBucket)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(r)
		if err != nil {
			return err
		}
		return b.Put([]byte(r.ID), encoded)
	})
	return err

}

func (s *BoltStore) GetResult(ID string) (Result, error) {
	val, err := s.getValue(resultsBucket, ID)

	res := Result{}

	if err != nil {
		return res, err
	}

	err = json.Unmarshal(val, &res)

	return res, err
}

func (s *BoltStore) GetMask(id string) ([]image.Rectangle, error) {
	key := id
	serialized, err := s.getValue(masksBucket, key)
	if err != nil {
		if err == KeyNotFoundError {
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

func (s *BoltStore) StoreMask(masks []image.Rectangle) (string, error) {
	key := RandStringBytes(10)

	serialized, err := json.Marshal(masks)
	if err != nil {
		return "", err
	}

	err = s.storeValue(masksBucket, key, serialized)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *BoltStore) GetImage(imgID string) (image.Image, error) {
	val, err := s.getValue(imagesBucket, imgID)

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

	err = s.storeValue(imagesBucket, imgID, buffer.Bytes())

	return imgID, err
}

func (s *BoltStore) GetBaseImageID(projectID, branch, target, browser string) (string, error) {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	return s.getStringValue(baseImagesBucket, key)
}

func (s *BoltStore) SetBaseImageID(baseImageID, projectID, branch, target, browser string) error {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	return s.storeStringValue(baseImagesBucket, key, baseImageID)
}

func (s *BoltStore) GetBaseMaskID(projectID, branch, target, browser string) (string, error) {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	return s.getStringValue(baseMasksBucket, key)
}

func (s *BoltStore) SetBaseMaskID(baseMaskID, projectID, branch, target, browser string) error {
	key := s.generateUniqueKey(projectID, branch, target, browser)
	return s.storeStringValue(baseMasksBucket, key, baseMaskID)
}
