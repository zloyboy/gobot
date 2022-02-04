package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/zloyboy/gobot/user"

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

var scripts = [4]string{"data/create_tbl_user.sql", "data/create_tbl_ill.sql", "data/create_tbl_vac.sql",
	"data/create_idx_name.sql"}

func createDb(db *sql.DB) error {
	for _, script := range scripts {
		if !runScript(db, script) {
			return fmt.Errorf("couldn't run %s ", script)
		}
	}
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
	ill []user.UserIll,
	countVac int,
	vac []user.UserVac) error {

	tx, _ := dbase.db.Begin()
	defer tx.Rollback()

	stmt, _ := tx.Prepare("INSERT INTO user" +
		"(teleId, created, modified, name, country, birth, gender, education, vaccine, origin, countIll, countVac)" +
		"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	stmt.Exec(teleId, date, date, name, country, birth, gender, education, vaccine, origin, countIll, countVac)

	for i := 0; i < countIll; i++ {
		stmt, _ = tx.Prepare("INSERT INTO userIllness (id, created, teleId, year, month, sign, degree) values(?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(nil, date, teleId, ill[i].Year, ill[i].Month, ill[i].Sign, ill[i].Degree)
	}

	for i := 0; i < countVac; i++ {
		stmt, _ = tx.Prepare("INSERT INTO userVaccine (id, created, teleId, year, month, kind, effect) values(?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(nil, date, teleId, vac[i].Year, vac[i].Month, vac[i].Kind, vac[i].Effect)
	}

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
