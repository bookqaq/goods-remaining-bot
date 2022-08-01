package database

import (
	"database/sql"
	"log"
	"sync"
)

var RecordSpace struct {
	InsertOne,
	UpdateType,
	DeleteOne,
	SelectByOwner,
	Auth *sql.Stmt
}

func recordSpaceOperations(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	tmp, err := db.Prepare("INSERT INTO recordSpace(owner, name, type) VALUES (?, ?, ?);")
	if err != nil {
		log.Panic(err)
	}
	RecordSpace.InsertOne = tmp

	tmp, err = db.Prepare("UPDATE recordSpace SET type=? WHERE name=? AND owner=?;")
	if err != nil {
		log.Panic(err)
	}
	RecordSpace.UpdateType = tmp

	tmp, err = db.Prepare("DELETE FROM recordSpace WHERE name=? AND owner=?;")
	if err != nil {
		log.Panic(err)
	}
	RecordSpace.DeleteOne = tmp

	tmp, err = db.Prepare("SELECT owner, name, type FROM recordSpace WHERE owner=?;")
	if err != nil {
		log.Panic(err)
	}
	RecordSpace.SelectByOwner = tmp

	tmp, err = db.Prepare("SELECT id FROM recordSpace WHERE name=? AND owner=?")
	if err != nil {
		log.Panic(err)
	}
	RecordSpace.Auth = tmp
}
