package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DB *sqlx.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		log.Print("error connecting to database", err)
	}
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("failed to read schema file: %v", err)
	}
	DB.Exec(string(schema))
}
