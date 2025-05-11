package main

import (
	"log"
	grpc "serverCalc/gprc"
	application "serverCalc/internal"
	"serverCalc/storage"
)

func main() {
	app := application.NewApplication()
	grpc.StartGRPCServer()
	app.RunServer()

	db, err := storage.InitDB("calc.db")
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()
}
