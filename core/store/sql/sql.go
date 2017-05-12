package sql

import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/theopticians/optician-api/core/store"
	"github.com/theopticians/optician-api/core/structs"
)

var schema = `
	CREATE TABLE masks (
			id integer,
			mask text,
	);

	CREATE TABLE results (
		id string,
		project string,
		branch string,
		batch string,
		target string,
		browser string,
		mask string,
		diffscore float,
		image string,
		baseimage string,
		diffimage string,
		diffclusters string,
		timestamp date
	);
`

const imagesPath = "images-storage"

type SqlStore struct {
	conn *sqlx.DB
}

func NewSqlStore(driver, url string) store.Store {
	conn, err := sqlx.Connect(driver, url)
	if err != nil {
		log.Fatal(err)
	}
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
	err := s.conn.Select(&results, "SELECT * FROM results WHERE batch=? ORDER BY timestamp DESC", batch)

	return results, err
}

func (s *SqlStore) GetBatchs() ([]structs.BatchInfo, error) {
	//err = db.Select(&names, "SELECT name FROM place LIMIT 10")
	//TODO
	batches := []structs.BatchInfo{}
	err := s.conn.Select(&batches, `"SELECT t1.ID, t1.Name, t1.Price, t1.Date 
	FROM   temp t1 
	INNER JOIN 
	(
		SELECT Max(date) date, name
		FROM   temp 
		GROUP BY name 
	) AS t2 
	ON t1.name = t2.name
	AND t1.date = t2.date 
	ORDER BY date DESC `)

	return batches, err

	return []structs.BatchInfo{}, nil
}

func (s *SqlStore) GetLastResult(projectID, branch, target, browser string) (structs.Result, error) {
	result := structs.Result{}
	err := s.conn.Get(&result, "SELECT * FROM results ORDER BY timestamp DESC LIMIT 1")

	return result, err
}

func (s *SqlStore) StoreResult(r structs.Result) error {
	_, err := s.conn.NamedExec("INSERT INTO results (id,project,branch,batch,target,browser,mask,diffscore,image,baseimage,diffimage,diffclusters,timestamp) VALUES (:id,:project,:branch,:batch,:target,:browser,:mask,:diffscore,:image,:baseimage,:diffimage,:diffclusters,:timestamp)", r)

	return err
}

func (s *SqlStore) GetResult(ID string) (structs.Result, error) {
	result := structs.Result{}
	err := s.conn.Get(&result, "SELECT * FROM results WHERE id=?", ID)

	return result, err
}

func (s *SqlStore) GetMask(id string) (structs.Mask, error) {
	mask := structs.Mask{}
	err := s.conn.Get(&mask, "SELECT mask FROM masks WHERE id=?", id)

	return mask, err
}

func (s *SqlStore) StoreMask(masks structs.Mask) (string, error) {
}

func imageAbsPath(id string) string {
	path, err := filepath.Abs(path.Join(imagesPath, id+".png"))
	if err != nil {
		panic("error generating imageAbsPath")
	}
	return path
}

// The image methods make sense, but in the SQL case we are encoding/decoding 2 times and we dont need to (API is png and DB is png)
func (s *SqlStore) GetImage(imgID string) (image.Image, error) {
	if _, err := os.Stat(imageAbsPath(imgID)); os.IsNotExist(err) {
		return nil, store.NotFoundError
	}

	reader, err := os.Open(imageAbsPath(imgID))

	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.New("Could not decode image")
	}

	return m, nil
}

func (s *SqlStore) StoreImage(img image.Image) (string, error) {
	var id string

	for {
		id = RandStringBytes(10)
		if _, err := os.Stat(imageAbsPath(id)); os.IsNotExist(err) {
			break
		}
	}

	out, err := os.Create(imageAbsPath(id))
	if err != nil {
		return "", err
	}

	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *SqlStore) GetBaseImageID(projectID, branch, target, browser string) (string, error) {
	var imgID string
	err := s.conn.Get(&imgID, "SELECT imageid FROM base_images WHERE projectid=? AND branch=? AND target=? AND browser=?", projectID, branch, target, browser)

	return imgID, err
}

func (s *SqlStore) SetBaseImageID(baseImageID, projectID, branch, target, browser string) error {
}

func (s *SqlStore) GetBaseMaskID(projectID, branch, target, browser string) (string, error) {
	var maskID string
	err := s.conn.Get(&maskID, "SELECT maskid FROM base_masks WHERE projectid=? AND branch=? AND target=? AND browser=?", projectID, branch, target, browser)

	return maskID, err
}

func (s *SqlStore) SetBaseMaskID(baseMaskID, projectID, branch, target, browser string) error {
}
