package database

import (
	"database/sql"
	"log"
	"sync"
)

var ImageStore struct {
	InsertOne,
	SelectByRS,
	DeleteOne,
	DeleteByRS *sql.Stmt
}

func imageStoreOperations(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	tmp, err := db.Prepare("INSERT INTO imageStore(priv, rs, url) VALUES (?, ?, ?);")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.InsertOne = tmp

	tmp, err = db.Prepare("SELECT url, name FROM imageStore WHERE rs=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.SelectByRS = tmp

	tmp, err = db.Prepare("DELETE FROM imageStore WHERE priv=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.DeleteOne = tmp

	tmp, err = db.Prepare("DELETE FROM imageStore WHERE rs=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.DeleteByRS = tmp
}