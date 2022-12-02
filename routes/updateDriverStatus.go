package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type updateStatusBody struct {
	DriverID   string `json:"driver_id"`
	RideID     string `json:"ride_id"`
	DriverLat  string `json:"latitude"`
	DriverLong string `json:"longitude"`
	Eta        int    `json:"eta"`
	HasRider   bool   `json:has_rider`
}

func updateDriverStatus(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/updateDriverStatus",
		Handler: func(c echo.Context) error {
			var body updateStatusBody

			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)
			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", body.DriverID)

			if driver != nil {
				driver.SetDataValue("current_latitude", body.DriverLat)
				driver.SetDataValue("current_longitude", body.DriverLong)
				driver.SetDataValue("has_rider", body.HasRider)

				app.Dao().SaveRecord(driver)
				if ride == nil {
					return c.NoContent(200)
				}
			}

			if ride != nil {
				// Update the ride with the new information
				ride.SetDataValue("eta", body.Eta)
				app.Dao().SaveRecord(ride)

				return c.NoContent(200)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
