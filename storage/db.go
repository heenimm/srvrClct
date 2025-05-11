package storage

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func InitDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		return nil, err
	}

	return db, nil
}

func ApplySchema(db *sqlx.DB, schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Printf("failed to read schema file: %v", err)
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Printf("failed to execute schema: %v", err)
		return err
	}
	return nil
}
