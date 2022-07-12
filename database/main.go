package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize() {
	db := conncet()

	var wg *sync.WaitGroup
	wg.Add(4)

	go imageStoreOperations(db, wg)
	go recordSpaceOperations(db, wg)
	go rsGroupMappingOperations(db, wg)
	go rsUserMappingOperations(db, wg)

	wg.Wait()
}

func conncet() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/db/sqlite3.db")
	if err != nil {
		log.Panic(err)
	}
	return db
}
