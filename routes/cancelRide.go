package routes

import (
	"net/http"
	"sync"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type cancelRideBody struct {
	RideID string `json:"ride_id"`
}

func cancelRide(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue, mutex *sync.RWMutex) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/cancelRide",
		Handler: func(c echo.Context) error {
			var body cancelRideBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)

			if ride != nil {
				if ride.GetBoolDataValue("needs_ride") && !ride.GetBoolDataValue("in_ride") {
					eventCode := ride.GetStringDataValue("event_id")
					app.Dao().DeleteRecord(ride)

					mutex.Lock()
					helpers.UpdateDriverQueues(app, eventCode, queues, nil, nil)
					mutex.Unlock()

					return c.NoContent(200)
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
