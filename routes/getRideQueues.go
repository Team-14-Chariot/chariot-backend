package routes

import (
	"net/http"

	. "github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getRideQueuesBody struct {
	EventID string `json:"event_id"`
}

type getRideQueuesResp struct {
	Queues []queueResp `json:"queues"`
}

type queueResp struct {
	DriverName string `json:"name"`
	Rides      []Ride `json:"rides"`
}

func getRideQueues(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getRideQueues",
		Handler: func(c echo.Context) error {
			var body getRideQueuesBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)
			if event != nil {
				resp := getRideQueuesResp{[]queueResp{}}

				drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
				driversRecords := GetEventDrivers(app, drivers_col, body.EventID)

				for _, driver := range driversRecords {
					rides, valid := queues[driver.Id]
					if valid {
						resp.Queues = append(resp.Queues, queueResp{DriverName: driver.GetStringDataValue("name"), Rides: rides.GetRides()})
					} else {
						resp.Queues = append(resp.Queues, queueResp{DriverName: driver.GetStringDataValue("name"), Rides: []Ride{}})
					}
				}

				return c.JSON(200, resp)
			} else {
				return c.NoContent(400)
			}
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
