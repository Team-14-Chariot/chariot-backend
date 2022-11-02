package helpers

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func GetNeededRides(app core.App, rides_col *models.Collection, event_id string) []models.Record {
	rides := GetAllRecords(app, rides_col)

	n := 0
	for _, ride := range rides {
		if (ride.GetDataValue("event_id") == event_id) && (ride.GetDataValue("needs_ride") == true) {
			rides[n] = ride
			n++
		}
	}

	rides = rides[:n]
	return rides
}
