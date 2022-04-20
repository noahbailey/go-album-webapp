package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Check if the file exists. Do not create it
func CheckDataFile(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// Create the database connection
func InitDB(filename string) (*sql.DB, error) {
	//open DB connection
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

//Create the table if needed
func CreateTable(db *sql.DB) {
	tbl := `
	CREATE TABLE IF NOT EXISTS album(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title TEXT NOT NULL,
		Artist TEXT NOT NULL,
		Price FLOAT
	);`

	log.Println("Creating table...")
	statement, err := db.Prepare(tbl)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Created table successfully.")
}
