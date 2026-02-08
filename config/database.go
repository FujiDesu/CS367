package config

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "students.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id TEXT PRIMARY KEY,
		name TEXT,
		major TEXT,
		gpa REAL
	)`)
	return db
}
