package main

import (
	application "serverCalc/internal"
	"serverCalc/storage"
)

func main() {
	app := application.NewApplication()
	storage.InitDB("calc.db")
	app.RunServer()
}
