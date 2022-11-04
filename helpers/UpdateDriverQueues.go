package helpers

import (
	"fmt"

	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/pocketbase/pocketbase"
)

func UpdateDriverQueues(app *pocketbase.PocketBase, eventID string, queues map[string]DriverQueue) {
	drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
	driversRecords := GetEventDrivers(app, drivers_col, eventID)
	rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
	ridesRecords := GetNeededRides(app, rides_col, eventID)
	drivers := ConvertToDriverObject(app, driversRecords)
	rides := ConvertToRideObject(ridesRecords)

	// for _, driver := range drivers {
	// 	for _, ride := range rides {
	// 		//g.AddEdge(driver.ID, ride.ID, graph.EdgeWeight(24))
	// 	}
	// }

	// for i, _ := range rides {
	// 	for j, _ := range rides {
	// 		if j > i {
	// 			// g.AddEdge(rides[i].ID, rides[j].ID, graph.EdgeWeight(12))
	// 			// g.AddEdge(rides[j].ID, rides[i].ID, graph.EdgeWeight(5))
	// 		}
	// 	}
	// }

	fmt.Println(drivers)
	fmt.Println(rides)
}
