package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var GoodsImages struct {
	InsertOne,
	SelectAll,
	SelectByName,
	DeleteOne,
	DeleteByName *sql.Stmt
}

var GroupUserMapping struct {
	InsertOne,
	InsertMany,
	SelectByOwner,
	DeleteOne,
	DeleteMany,
	Exist *sql.Stmt
}

func Initialize() {
	db := conncet()
	goodsOperation(db)
}

func conncet() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/db/sqlite3.db")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func goodsOperation(db *sql.DB) {
	tmp, err := db.Prepare("INSERT INTO remainingGoodsImages(name, url) values (?,?);")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.InsertOne = tmp

	tmp, err = db.Prepare("SELECT priv, name, url FROM remainingGoodsImages;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.SelectAll = tmp

	tmp, err = db.Prepare("SELECT priv, name, url FROM remainingGoodsImages where name=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.SelectByName = tmp

	tmp, err = db.Prepare("DELETE FROM remainingGoodsImages WHERE priv=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.DeleteOne = tmp

	tmp, err = db.Prepare("DELETE FROM remainingGoodsImages WHERE name=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.DeleteByName = tmp
}

func groupUserMappingOperation(db *sql.DB) {
	tmp, err := db.Prepare("SELECT 1 FROM groupUserMappings WHERE id=? LIMIT 1")
	if err != nil {
		log.Panic(err)
	}
	GroupUserMapping.Exist = tmp
}
