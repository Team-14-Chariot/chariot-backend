package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
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

func getRide(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getRide",
		Handler: func(c echo.Context) error {
			var body getRideBody

			c.Bind(&body)

			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", body.DriverID)

			if driver != nil {
				rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
				rides := helpers.GetAllRecords(app, rides_col)

				if len(rides) > 0 {
					ride := rides[0]

					return c.JSON(200, map[string]interface{}{
						"ride_id":          ride.Id,
						"source_latitude":  ride.GetDataValue("origin_latitude"),
						"source_longitude": ride.GetDataValue("origin_longitude"),
						"dest_latitude":    ride.GetDataValue("dest_latitude"),
						"dest_longitude":   ride.GetDataValue("dest_longitude"),
						"rider_name":       "Jeff",
					})
				}
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
