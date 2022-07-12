package database

import (
	"database/sql"
	"log"
	"sync"
)

var RSGroupMapping struct {
	InsertOne,
	SelectByRS,
	DeleteOne *sql.Stmt
}

func rsGroupMappingOperations(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	tmp, err := db.Prepare("INSERT INTO rsGroupMapping(rs, gp) VALUES(?, ?);")
	if err != nil {
		log.Panic(err)
	}
	RSGroupMapping.InsertOne = tmp

	tmp, err = db.Prepare("SELECT gp FROM rsGroupMapping WHERE rs=?;")
	if err != nil {
		log.Panic(err)
	}
	RSGroupMapping.SelectByRS = tmp

	tmp, err = db.Prepare("DELETE FROM rsGroupMapping WHERE rs=? AND gp=?;")
	if err != nil {
		log.Panic(err)
	}
	RSGroupMapping.DeleteOne = tmp

}
