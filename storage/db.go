package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func InitDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		return nil, err
	}

	schema, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		log.Printf("failed to read schema file: %v", err)
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Printf("failed to execute schema: %v", err)
		return nil, err
	}

	return db, nil
}
