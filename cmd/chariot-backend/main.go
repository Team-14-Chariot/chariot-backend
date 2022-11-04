package main

import (
	"log"

	. "github.com/Team-14-Chariot/chariot-backend/models"
	"github.com/Team-14-Chariot/chariot-backend/routes"
	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	queues := make(map[string]DriverQueue)

	routes.Routes(app, queues)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
