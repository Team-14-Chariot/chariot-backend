package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type endRideBody struct {
	RideID string `json:"ride_id"`
}

func endRide(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/endRide",
		Handler: func(c echo.Context) error {
			var body endRideBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)
			ride.SetDataValue("in_ride", false)
			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", ride.GetStringDataValue("driver_id"))
			driver.SetDataValue("has_rider", false)
			driver.SetDataValue("in_ride", false)

			app.Dao().SaveRecord(ride)
			app.Dao().SaveRecord(driver)
			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
