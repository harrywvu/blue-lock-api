package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB){
	db, err := sql.Open("sqlite3", "./blue-lock.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA foreign_keys = ON")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS players(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age	INTEGER NOT NULL,
		nelTeam TEXT , 
		primaryPosition TEXT NOT NULL,
		currentBlueLockRank INTEGER
	)`)
	if err != nil {log.Fatalf("Failed to create table: %v", err)}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stats(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		playerid INTEGER NOT NULL,
		overall REAL NOT NULL,
		offense REAL NOT NULL,
		shooting REAL NOT NULL,
		speed REAL NOT NULL,
		defense REAL NOT NULL,
		passing REAL NOT NULL,
		dribbling REAL NOT NULL,
		FOREIGN KEY (playerid) REFERENCES players(id)
	)`)
	if err != nil {log.Fatalf("Failed to create table: %v", err)}

	return db
}