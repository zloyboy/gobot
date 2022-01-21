package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Init() error {
	sqliteDatabase, err := sql.Open("sqlite3", "data/stat.db")
	if err != nil {
		return err
	}
	defer sqliteDatabase.Close()

	return nil
}
