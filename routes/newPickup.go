package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type requestRideBody struct {
	EventID    string `json:"event_id"`
	OriginLat  string `json:"origin_latitude"`
	OriginLong string `json:"origin_longitude"`
	DestLat    string `json:"dest_latitude"`
	DestLong   string `json:"dest_longitude"`
	RiderName  string `json:"rider_name"`
	GroupSize  int    `json:"group_size"`
	Ride_id    string `json:"ride_id"`
}

func requestRide(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/requestRide",
		Handler: func(c echo.Context) error {
			var body requestRideBody
			c.Bind(&body)

			drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
			drivers := helpers.GetEventDrivers(app, drivers_col, body.EventID)

			if len(drivers) > 0 {
				rides, _ := app.Dao().FindCollectionByNameOrId("rides")
				rider, _ := app.Dao().FindFirstRecordByData(rides, "ride_id", body.Ride_id)

				
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
