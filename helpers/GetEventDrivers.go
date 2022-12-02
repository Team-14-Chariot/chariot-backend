package helpers

import (
	. "github.com/Team-14-Chariot/chariot-backend/models"

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

func GetAllEventDrivers(app *pocketbase.PocketBase, drivers_col *models.Collection, event_id string) []models.Record {
	drivers := GetAllRecords(app, drivers_col)

	n := 0
	for _, driver := range drivers {
		if driver.GetDataValue("event_id") == event_id {
			drivers[n] = driver
			n++
		}
	}

	drivers = drivers[:n]
	return drivers
}

func ConvertToDriverObject(app *pocketbase.PocketBase, drivers []models.Record) []*Driver {
	var driverObjects []*Driver

	for _, driver := range drivers {
		if driver.GetBoolDataValue("in_ride") {
			rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride := GetDriversCurrentRide(app, rides_col, driver.Id)

			etaToLocation := ride.GetFloatDataValue("eta")

			if !driver.GetBoolDataValue("has_rider") {
				etaToLocation += ride.GetFloatDataValue("ride_length")
			}

			driverObjects = append(driverObjects, &Driver{
				ID:            driver.Id,
				Capacity:      driver.GetIntDataValue("car_capacity"),
				CurrentLat:    ride.GetStringDataValue("dest_latitude"),
				CurrentLong:   ride.GetStringDataValue("dest_longitude"),
				EtaToLocation: etaToLocation,
			})
		} else {
			driverObjects = append(driverObjects, &Driver{
				ID:            driver.Id,
				Capacity:      driver.GetIntDataValue("car_capacity"),
				CurrentLat:    driver.GetStringDataValue("current_latitude"),
				CurrentLong:   driver.GetStringDataValue("current_longitude"),
				EtaToLocation: 0,
			})
		}
	}

	return driverObjects
}
