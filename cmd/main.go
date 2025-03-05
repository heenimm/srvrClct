package main

import (
	application "serverCalc/internal"
)

func main() {
	app := application.NewApplication()
	app.RunServer()
}
