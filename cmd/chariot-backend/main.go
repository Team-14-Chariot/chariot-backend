package main

import (
	"log"

	"github.com/Team-14-Chariot/chariot-backend/routes"
	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	routes.Routes(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
