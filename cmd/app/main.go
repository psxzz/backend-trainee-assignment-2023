package main

import (
	"log"

	"github.com/psxzz/backend-trainee-assignment/pkg/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
