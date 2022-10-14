package helpers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func GetEventDrivers(app *pocketbase.PocketBase, drivers_col *models.Collection, event_id string) []models.Record {
	drivers := GetAllRecords(app, drivers_col)

	n := 0
	for _, driver := range drivers {
		if (driver.GetDataValue("event_id") == event_id) && (driver.GetDataValue("active") == true) {
			drivers[n] = driver
			n++
		}
	}

	drivers = drivers[:n]
	return drivers
}
