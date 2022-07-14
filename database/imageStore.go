package database

import (
	"database/sql"
	"log"
	"sync"
)

var ImageStore struct {
	InsertOne,
	SelectOne,
	SelectByRS,
	UpdateOne,
	DeleteOne,
	DeleteByRS,
	Exist *sql.Stmt
}

func imageStoreOperations(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	tmp, err := db.Prepare("INSERT INTO imageStore(priv, rs, url) VALUES (?, ?, ?);")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.InsertOne = tmp

	tmp, err = db.Prepare("SELECT priv, url, name FROM imageStore WHERE priv=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.SelectOne = tmp

	tmp, err = db.Prepare("SELECT priv, url, name FROM imageStore WHERE rs=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.SelectByRS = tmp

	tmp, err = db.Prepare("UPDATE imageStore SET url=? WHERE priv=?;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.UpdateOne = tmp

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

	tmp, err = db.Prepare("SELECT priv FROM imageStore WHERE priv=? LIMIT 1;")
	if err != nil {
		log.Panic(err)
	}
	ImageStore.Exist = tmp
}
