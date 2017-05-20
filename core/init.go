package core

import (
	"fmt"
	"os"

	"github.com/theopticians/optician-api/core/store"
	"github.com/theopticians/optician-api/core/store/bolt"
	"github.com/theopticians/optician-api/core/store/sql"
)

var db store.Store

func init() {
	storeType := os.Getenv("STORE_TYPE")
	if storeType == "sql" {
		println("Using backend: sql")
		sqlHost := os.Getenv("SQL_HOST")
		sqlPort := os.Getenv("SQL_PORT")
		db = sql.NewSqlStore("postgres", fmt.Sprintf("postgresql://root@%s:%s/optician?sslmode=disable", sqlHost, sqlPort))
	} else {
		println("Using backend: boltdb")
		db = bolt.NewBoltStore("./optician.db")
	}
}
