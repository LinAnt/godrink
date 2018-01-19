package database

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Card represents the full data we store per card
type Card struct {
	ID         string    `db:"card_id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	MatNbr     int64     `db:"mat_nbr"`
	Credit     int64     `db:"credit"`
	LastUpdate time.Time `db:"last_update"`
	LastTopUp  time.Time `db:"last_top_up"`
}

// CreateDatabase sets you up with a new godrink compatible DB
func CreateDatabase(path string) error {

	if _, err := os.Stat(path); err == nil {
		return errors.New("new db should not overwrite file")
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE Cards (
		card_id TEXT UNIQUE NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		mat_nbr INT64 NOT NULL,
		last_update DATETIME NOT NULL,
		last_top_up DATETIME NOT NULL,
		credit INT64 NOT NULL CHECK(credits > 0)
		);
	`
	_, err = db.Exec(sqlStmt)
	return err
}

// Database represents our sqlx instance
type Database struct {
	*sqlx.DB
}

//GetDB returns new db instance, remember to close afterwards
func GetDB(filepath string) (Database, error) {
	db, err := sqlx.Open("sqlite3", filepath)
	return Database{db}, err
}

// GetCard returns card of cardID
func (db *Database) GetCard(cardID string) (Card, error) {
	card := Card{}
	return card, db.Select(&card, "SELECT * FROM Cards WHERE card_id = $1", cardID)
}

// GetCards returns all cards in the Database
func (db *Database) GetCards() ([]Card, error) {
	cards := []Card{}
	return cards, db.Select(&cards, "SELECT * FROM Cards")
}

// GetBalance returns current balance for the specified card_id
func (db *Database) GetBalance(cardID string) (int64, error) {
	var balance int64
	return balance, db.Select(balance, "SELECT credit FROM Cards WHERE card_id = $1", cardID)
}

// ChangeBalance changes balance of cardID by amount
func (db *Database) ChangeBalance(cardID string, amount int64) error {
	curBalance, err := db.GetBalance(cardID)
	if err != nil {
		return err
	}
	return db.setBalance(cardID, curBalance+amount)
}

// NewUser adds a new user to the database
func (db *Database) NewUser(cardID, name, surname string, matNbr, amount int64) error {
	tx := db.MustBegin()
	// TODO: Don't drink and code :P
	_, err := tx.NamedExec(`INSERT INTO Cards (card_id, name, surname, mat_nbr, last_update, last_top_up, credit) 
	VALUES (:card_id, :name, :surname, :mat_nbr, :last_update, :last_top_up, :credit)`, &Card{
		ID:         cardID,
		Name:       name,
		Surname:    surname,
		MatNbr:     matNbr,
		LastTopUp:  time.Now(),
		LastUpdate: time.Now(),
		Credit:     amount,
	})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) setBalance(cardID string, amount int64) error {

	query := `UPDATE table_name
	SET credit = ?
	WHERE card_id = ?;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(amount, cardID)
	return err
}
