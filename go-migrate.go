package main

import (
	"github.com/akaumov/go-migrate/lib"
)

func main() {
	app := lib.NewApp("0.0.2", &lib.EmptyHandler{})
	app.Execute()
}