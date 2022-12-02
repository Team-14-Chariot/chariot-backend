package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getEtaBody struct {
	RideID string `json:"ride_id"`
}

func getEta(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getEta",
		Handler: func(c echo.Context) error {
			var body getEtaBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)

			if ride != nil {
				if ride.GetBoolDataValue("in_ride") {
					return c.JSON(200, map[string]interface{}{"eta": ride.GetDataValue("eta")})
				}

				driver := helpers.UpdateDriverQueues(app, ride.GetStringDataValue("event_id"), queues, nil, &Ride{ID: body.RideID})
				driverQueue, valid := queues[driver.ID]
				if valid {
					eta := helpers.CalculateTotalTripLength(driverQueue, driver, body.RideID)

					return c.JSON(200, map[string]interface{}{"eta": eta})
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
