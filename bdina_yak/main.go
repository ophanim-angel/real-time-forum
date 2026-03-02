package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	dsn := "./database/forum.db?_foreign_keys=on&_journal_mode=WAL&_synchronous=normal"
	db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database not reachable: ", err)
	}

	db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)
	// db.SetConnMaxLifetime(0)

	if err = CreateTables(); err != nil {
		log.Fatal("Problem while creating tables ", err)
	} else {
		log.Println("Database connected successfully and configured.")
	}
}
