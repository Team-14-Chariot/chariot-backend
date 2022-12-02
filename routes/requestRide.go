package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"
	. "github.com/Team-14-Chariot/chariot-backend/tools"

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
}

func requestRide(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
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

				events, _ := app.Dao().FindCollectionByNameOrId("events")
				event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)

				if event.GetDataValue("accept_rides") == true {
					newRide := models.NewRecord(rides)

					newRide.SetDataValue("event_id", body.EventID)
					newRide.SetDataValue("origin_latitude", body.OriginLat)
					newRide.SetDataValue("origin_longitude", body.OriginLong)
					newRide.SetDataValue("dest_latitude", body.DestLat)
					newRide.SetDataValue("dest_longitude", body.DestLong)
					newRide.SetDataValue("rider_name", body.RiderName)
					newRide.SetDataValue("needs_ride", true)
					newRide.SetDataValue("in_ride", false)
					newRide.SetDataValue("group_size", body.GroupSize)
					length := CalculateRideLength(*newRide)
					newRide.SetDataValue("ride_length", length)

					app.Dao().SaveRecord(newRide)

					helpers.UpdateDriverQueues(app, body.EventID, queues, nil, nil)

					return c.JSON(200, map[string]interface{}{"ride_id": newRide.Id})
				}

				return c.NoContent(512)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
