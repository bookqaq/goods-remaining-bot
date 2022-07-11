package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var GoodsImages struct {
	InsertOne    *sql.Stmt
	SelectAll    *sql.Stmt
	SelectByName *sql.Stmt
	DeleteOne    *sql.Stmt
	DeleteByName *sql.Stmt
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
	tmp, err := db.Prepare("INSERT INTO goodsImages(name, url) values (?,?);")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.InsertOne = tmp

	tmp, err = db.Prepare("SELECT priv, name, url FROM goodsImages;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.SelectAll = tmp

	tmp, err = db.Prepare("SELECT priv, name, url FROM goodsImages where name=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.SelectByName = tmp

	tmp, err = db.Prepare("DELETE FROM goodsImages WHERE priv=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.DeleteOne = tmp

	tmp, err = db.Prepare("DELETE FROM goodsImages where name=?;")
	if err != nil {
		log.Panic(err)
	}
	GoodsImages.DeleteByName = tmp
}
