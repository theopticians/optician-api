package sql

import (
	"bytes"
	"database/sql"
	"errors"
	"image"
	"image/png"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/theopticians/optician-api/core/store"
	"github.com/theopticians/optician-api/core/structs"

	// Import postgres driver.
	_ "github.com/lib/pq"
)

var schema = `
	CREATE TABLE IF NOT EXISTS masks (
		id string,
		mask text,
		PRIMARY KEY( id )
	);

	CREATE TABLE IF NOT EXISTS images (
		id string,
		image bytea,
		PRIMARY KEY( id )
	);

	CREATE TABLE IF NOT EXISTS results (
		id string,
		project string,
		branch string,
		batch string,
		target string,
		browser string,
		maskid string,
		diffscore float,
		imageid string,
		baseimageid string,
		diffimageid string,
		diffclusters string,
		timestamp date,
		PRIMARY KEY( id ),
		CONSTRAINT UQ_result UNIQUE ( project, branch, batch, target, browser )
	);

	CREATE TABLE IF NOT EXISTS base_images (
		project string,
		branch string,
		target string,
		browser string,
		imageid string,
		PRIMARY KEY( project, branch,target, browser )
	);

	CREATE TABLE IF NOT EXISTS base_masks (
		project string,
		branch string,
		target string,
		browser string,
		maskid string,
		PRIMARY KEY( project, branch,target, browser )
	);
`

type SqlStore struct {
	conn *sqlx.DB
}

func NewSqlStore(driver, url string) store.Store {
	conn, err := sqlx.Connect(driver, url)
	if err != nil {
		log.Fatal(err)
	}

	conn.MustExec("CREATE DATABASE IF NOT EXISTS optician")
	conn.MustExec("SET DATABASE = optician;")
	conn.MustExec(schema)

	return &SqlStore{conn}
}

func (s *SqlStore) Close() {
	s.conn.Close()
}

func (s *SqlStore) GetResults() ([]structs.Result, error) {
	results := []structs.Result{}
	err := s.conn.Select(&results, "SELECT * FROM results ORDER BY timestamp DESC")

	return results, err
}

func (s *SqlStore) GetResultsByBatch(batch string) ([]structs.Result, error) {
	results := []structs.Result{}
	err := s.conn.Select(&results, "SELECT * FROM results WHERE batch=$1 ORDER BY timestamp DESC", batch)

	return results, err
}

func (s *SqlStore) GetBatchs() ([]structs.BatchInfo, error) {
	batches := []structs.BatchInfo{}
	err := s.conn.Select(&batches, `
	SELECT t1.batch AS id, t1.timestamp, t1.project, COALESCE(t3.failed,0) AS failed FROM results AS t1 JOIN (SELECT batch, max(timestamp) AS maxdate FROM results GROUP BY batch) AS t2 ON (t1.batch = t2.batch) AND (t1.timestamp = t2.maxdate) FULL JOIN (SELECT batch, count(*) AS failed FROM results WHERE diffscore>0 GROUP BY batch) AS t3 ON (t1.batch = t3.batch) GROUP BY t1.batch, t1.timestamp, t1.project, failed ORDER BY t1.timestamp DESC
`)

	return batches, err
}

func (s *SqlStore) GetLastResult(projectID, branch, target, browser string) (structs.Result, error) {
	result := structs.Result{}
	err := s.conn.Get(&result, "SELECT * FROM results ORDER BY timestamp DESC LIMIT 1")

	if err == sql.ErrNoRows {
		return result, store.NotFoundError
	}

	return result, err
}

func (s *SqlStore) StoreResult(r structs.Result) error {
	_, err := s.conn.NamedExec("INSERT INTO results (id,project,branch,batch,target,browser,maskid,diffscore,imageid,baseimageid,diffimageid,diffclusters,timestamp) VALUES (:id,:project,:branch,:batch,:target,:browser,:maskid,:diffscore,:imageid,:baseimageid,:diffimageid,:diffclusters,:timestamp)", r)

	return err
}

func (s *SqlStore) GetResult(ID string) (structs.Result, error) {
	result := structs.Result{}
	err := s.conn.Get(&result, "SELECT * FROM results WHERE id=$1", ID)

	if err == sql.ErrNoRows {
		return result, store.NotFoundError
	}

	return result, err
}

func (s *SqlStore) GetMask(id string) (structs.Mask, error) {
	mask := structs.Mask{}
	err := s.conn.Get(&mask, "SELECT mask FROM masks WHERE id=$1", id)

	if err == sql.ErrNoRows {
		return nil, store.NotFoundError
	}

	return mask, err
}

func (s *SqlStore) StoreMask(mask structs.Mask) (string, error) {
	id := RandStringBytes(10)
	s.conn.MustExec("INSERT INTO masks (id, mask) VALUES ($1, $2)", id, mask)
	return id, nil
}

// The image methods make sense, but in the SQL case we are encoding/decoding 2 times and we dont need to (API is png and DB is png)
func (s *SqlStore) GetImage(imgID string) (image.Image, error) {
	imageBytes := []byte{}
	err := s.conn.Get(&imageBytes, "SELECT image FROM images WHERE id=$1", imgID)

	if err == sql.ErrNoRows {
		return nil, store.NotFoundError
	}

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(imageBytes)
	m, _, err := image.Decode(buff)
	if err != nil {
		return nil, errors.New("Could not decode image")
	}

	return m, nil
}

func (s *SqlStore) StoreImage(img image.Image) (string, error) {
	id := RandStringBytes(10)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return "", err
	}
	imageBytes := buf.Bytes()

	s.conn.MustExec("INSERT INTO images (id, image) VALUES ($1, $2)", id, imageBytes)
	return id, nil
}

func (s *SqlStore) GetBaseImageID(projectID, branch, target, browser string) (string, error) {
	var imgID string
	err := s.conn.Get(&imgID, "SELECT imageid FROM base_images WHERE project=$1 AND branch=$2 AND target=$3 AND browser=$4", projectID, branch, target, browser)

	if err == sql.ErrNoRows {
		return "", store.NotFoundError
	}

	return imgID, err
}

func (s *SqlStore) SetBaseImageID(baseImageID, projectID, branch, target, browser string) error {
	_, err := s.GetBaseImageID(projectID, branch, target, browser)

	if err == store.NotFoundError {
		s.conn.MustExec("INSERT INTO base_images (project, branch, target, browser, imageid) VALUES($1, $2, $3, $4, $5)", projectID, branch, target, browser, baseImageID)
	} else if err == nil {
		s.conn.MustExec("UPDATE base_images SET imageid=$5 WHERE project=$1 AND branch=$2 AND target=$3 AND browser=$4", projectID, branch, target, browser, baseImageID)
	}

	return nil
}

func (s *SqlStore) GetBaseMaskID(projectID, branch, target, browser string) (string, error) {
	var maskID string
	err := s.conn.Get(&maskID, "SELECT maskid FROM base_masks WHERE project=$1 AND branch=$2 AND target=$3 AND browser=$4", projectID, branch, target, browser)

	if err == sql.ErrNoRows {
		return "", store.NotFoundError
	}

	return maskID, err
}

func (s *SqlStore) SetBaseMaskID(baseMaskID, projectID, branch, target, browser string) error {
	s.conn.MustExec("INSERT INTO base_masks (project, branch, target, browser, maskid) VALUES($1, $2, $3, $4, $5) ON DUPLICATE KEY UPDATE maskid=$5", projectID, branch, target, browser, baseMaskID)
	return nil
}
