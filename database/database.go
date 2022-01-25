package database

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Dbase struct {
	db *sql.DB
}

func createDb(db *sql.DB) error {
	// read scripts
	content, err := ioutil.ReadFile("data/createtbl.sql")
	if err != nil {
		log.Print("Script createtbl.sql is not found")
		return err
	}
	cursor, _ := db.Prepare(string(content))
	cursor.Exec()

	content, err = ioutil.ReadFile("data/createidx.sql")
	if err != nil {
		log.Print("Script createidx.sql is not found")
		return err
	}
	cursor, _ = db.Prepare(string(content))
	cursor.Exec()

	return nil
}

func InitDb() *Dbase {
	db, err := sql.Open("sqlite3", "data/stat.db")
	if err != nil || db == nil {
		log.Print("Can't init sqlite3")
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Print("Can't ping sqlite3")
		return nil
	}
	var tableNname string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user'").Scan(&tableNname)
	if err != nil {
		log.Print("Table is not found, trying to create DB")
		err = createDb(db)
	}
	if err != nil {
		return nil
	}

	return &Dbase{db: db}
}

func (dbase *Dbase) CheckIdName(id int64) bool {
	return false
}
