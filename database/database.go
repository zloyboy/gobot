package database

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Dbase struct {
	db *sql.DB
}

func runScript(db *sql.DB, path string) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Script %s is not found", path)
		return false
	}
	stmt, err := db.Prepare(string(content))
	if err == nil {
		stmt.Exec()
	}
	return err == nil
}

func createDb(db *sql.DB) error {
	if runScript(db, "data/create_tbl_user.sql") {
		if runScript(db, "data/create_tbl_ill.sql") {
			if runScript(db, "data/create_tbl_vac.sql") {
				if runScript(db, "data/create_idx_name.sql") {
					return nil
				}
			}
		}
	}
	return errors.New("createDb fails")
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
	err := dbase.db.QueryRow("SELECT name FROM user WHERE teleId=?", id).Scan(&userNname)
	return userNname, err
}

func (dbase *Dbase) Insert(
	teleId int64,
	date string,
	name string,
	country string,
	birth int,
	gender int,
	education string,
	origin string,
	vaccine string,
	countIll int,
	countVac int) error {

	tx, _ := dbase.db.Begin()
	defer tx.Rollback()

	stmt, _ := tx.Prepare("INSERT INTO user" +
		"(teleId, created, name, country, birth, gender, education, vaccine, origin, countIll, countVac)" +
		"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	stmt.Exec(teleId, date, name, country, birth, gender, education, vaccine, origin, countIll, countVac)

	stmt, _ = tx.Prepare("INSERT INTO userIllness (id, created, teleId, year, month) values(?, ?, ?, ?, ?)")
	stmt.Exec(nil, date, teleId, 2021, 1)

	stmt, _ = tx.Prepare("INSERT INTO userVaccine (id, created, teleId, year, month) values(?, ?, ?, ?, ?)")
	stmt.Exec(nil, date, teleId, 2022, 2)

	tx.Commit()

	return nil
}

func (dbase *Dbase) CountUsers() int {
	var count = 0
	dbase.db.QueryRow("SELECT count(*) FROM user").Scan(&count)
	return count
}

func (dbase *Dbase) CountIll() int {
	var res = 0
	dbase.db.QueryRow("SELECT count(*) FROM user WHERE 0<countIll").Scan(&res)
	return res
}
