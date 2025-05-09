package main

import (
	"log"
	"net/http"
	grpc "serverCalc/gprc"
	"serverCalc/handlers"
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

	http.HandleFunc("/api/v1/register", handlers.RegisterHandler(db))
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
