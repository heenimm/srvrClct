package main

import (
	"log"
	grpc "serverCalc/gprc"
	application "serverCalc/internal"
	"serverCalc/storage"
)

func main() {
	db, err := storage.InitDB("calc.db")
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	app := application.NewApplication(db)
	go app.RunServer()
	grpc.StartGRPCServer()

}
