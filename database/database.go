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

func (dbase *Dbase) CheckIdName(id int64) (string, error) {
	var userNname = ""
	err := dbase.db.QueryRow("SELECT name FROM user WHERE id=?", id).Scan(&userNname)
	return userNname, err
}

func (dbase *Dbase) NewId(id int64) bool {
	var userID = 0
	dbase.db.QueryRow("SELECT id FROM user WHERE id=?", id).Scan(&userID)
	return userID == 0
}

func (dbase *Dbase) Insert(id int64, date string, name string, age int, res int) error {
	cursor, _ := dbase.db.Prepare("INSERT INTO user(id, created, name, age, res) values(?, ?, ?, ?, ?)")
	cursor.Exec(id, date, name, age, res)
	return nil
}

func (dbase *Dbase) CountUsers() int {
	var count = 0
	dbase.db.QueryRow("SELECT count(*) FROM user").Scan(&count)
	return count
}

func (dbase *Dbase) CountIll() int {
	var res = 0
	dbase.db.QueryRow("SELECT count(*) FROM user WHERE res=1").Scan(&res)
	return res
}

func (dbase *Dbase) CountAgeGroup(age int) int {
	var count = 0
	dbase.db.QueryRow("SELECT count(*) FROM user WHERE age=?", age).Scan(&count)
	return count
}

func (dbase *Dbase) CountAgeGroupIll(age int) int {
	var res = 0
	dbase.db.QueryRow("SELECT count(*) FROM user WHERE age=? AND res=1", age).Scan(&res)
	return res
}
