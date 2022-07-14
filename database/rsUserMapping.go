package database

import (
	"database/sql"
	"log"
	"sync"
)

var RSUserMapping struct {
	InsertOne,
	SelectOne,
	SelectRSByQQ,
	SelectByRS,
	DeleteOne,
	DeleteByRS *sql.Stmt
}

func rsUserMappingOperations(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	tmp, err := db.Prepare("INSERT INTO rsUserMapping(rs, dst) VALUES(?, ?);")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.InsertOne = tmp

	tmp, err = db.Prepare("SELECT owner, name, type, qq FROM opUserGetRS WHERE name=? AND qq=? LIMIT 1;")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.SelectOne = tmp

	tmp, err = db.Prepare("SELECT owner, name, type, qq FROM opUserGetRS WHERE qq=?;")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.SelectRSByQQ = tmp

	tmp, err = db.Prepare("SELECT id, rs, dst FROM rsUserMapping WHERE rs=?;")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.SelectByRS = tmp

	tmp, err = db.Prepare("DELETE FROM rsUserMapping WHERE rs=? AND dst=?;")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.DeleteOne = tmp

	tmp, err = db.Prepare("DELETE FROM rsUserMapping WHERE rs=?")
	if err != nil {
		log.Panic(err)
	}
	RSUserMapping.DeleteByRS = tmp

}
