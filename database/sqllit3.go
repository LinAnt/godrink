package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// CreateDatabase sets you up with a new godrink compatible DB
func CreateDatabase(path string) error {
	if _, err := os.Stat(path); err == nil {
		return errors.New("file already exists...")
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table Cards (
		card_id TEXT UNIQUE NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		mat_nbr INTEGER NOT NULL,
		last_purchase DATETIME,
		last_top_up DATETIME NOT NULL,
		credits INT64 CHECK(credits > 0)
		);
	`
	_, err = db.Exec(sqlStmt)
	return err
}
