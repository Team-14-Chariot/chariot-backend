package routes

import (
	"fmt"
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getRideBody struct {
	DriverID    string `json:"driver_id"`
	CurrentLat  string `json:"current_latitude"`
	CurrentLong string `json:"current_longitude"`
}

func getRide(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getRide",
		Handler: func(c echo.Context) error {
			var body getRideBody

			c.Bind(&body)

			drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers_col, "id", body.DriverID)

			if driver != nil && driver.GetBoolDataValue("active") {
				driver.SetDataValue("current_latitude", body.CurrentLat)
				driver.SetDataValue("current_longitude", body.CurrentLong)
				driver.SetDataValue("has_rider", false)
				app.Dao().SaveRecord(driver)
				helpers.UpdateDriverQueues(app, driver.GetStringDataValue("event_id"), queues, nil)

				driverQueue, valid := queues[driver.Id]
				if valid {
					ride := driverQueue.PopRide()
					if ride != nil {
						rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
						rideRecord, _ := app.Dao().FindFirstRecordByData(rides_col, "id", ride.ID)

						rideRecord.SetDataValue("needs_ride", false)
						rideRecord.SetDataValue("in_ride", true)
						rideRecord.SetDataValue("driver_id", body.DriverID)
						app.Dao().SaveRecord(rideRecord)

						driver.SetDataValue("in_ride", true)
						app.Dao().SaveRecord(driver)

						return c.JSON(200, map[string]interface{}{
							"ride_id":          ride.ID,
							"source_latitude":  ride.OriginLat,
							"source_longitude": ride.OriginLong,
							"dest_latitude":    ride.DestLat,
							"dest_longitude":   ride.DestLong,
							"rider_name":       ride.Name,
						})
					} else {
						return c.NoContent(201)
					}
				} else {
					fmt.Println("Here2")
					return c.NoContent(400)
				}
			}

			fmt.Println("Here1")
			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
