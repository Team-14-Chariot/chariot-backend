package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"
	"github.com/Team-14-Chariot/chariot-backend/tools"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getRouteETABody struct {
	EventID    string `json:"event_id"`
	OriginLat  string `json:"origin_latitude"`
	OriginLong string `json:"origin_longitude"`
	DestLat    string `json:"dest_latitude"`
	DestLong   string `json:"dest_longitude"`
}

func getRouteETA(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getRouteETA",
		Handler: func(c echo.Context) error {
			var body getRouteETABody
			c.Bind(&body)

			drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
			drivers := helpers.GetEventDrivers(app, drivers_col, body.EventID)

			if len(drivers) > 0 {
				events, _ := app.Dao().FindCollectionByNameOrId("events")
				event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)

				if event != nil && event.GetDataValue("accept_rides") == true {
					rideLength := tools.CalculateLength(body.OriginLat, body.OriginLong, body.DestLat, body.DestLong)

					testRide := &Ride{
						ID:         "ETA_RIDE",
						OriginLat:  body.OriginLat,
						OriginLong: body.OriginLong,
						DestLat:    body.DestLat,
						DestLong:   body.DestLong,
						RideLength: rideLength,
					}

					testQueues := make(map[string]*DriverQueue)

					driver := helpers.UpdateDriverQueues(app, body.EventID, testQueues, testRide, nil)
					driverQueue, valid := testQueues[driver.ID]

					eta := 0.0
					if valid {
						eta = helpers.CalculateTotalTripLength(driverQueue, driver, testRide.ID)
					}

					return c.JSON(200, map[string]interface{}{"eta": eta})
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
