package routes

import (
	"net/http"

	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getRideQueueBody struct {
	DriverID string `json:"driver_id"`
}

type getRideQueueResp struct {
	Queue []rideQueueResp `json:"queue"`
}

type rideQueueResp struct {
	Name string `json:"name"`
}

func getRideQueue(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getRideQueue",
		Handler: func(c echo.Context) error {
			var body getRideQueueBody
			c.Bind(&body)

			drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers_col, "id", body.DriverID)

			if driver != nil {
				resp := getRideQueueResp{Queue: []rideQueueResp{}}
				queue, valid := queues[driver.Id]

				if valid {
					rides := queue.GetRides()

					for _, ride := range rides {
						resp.Queue = append(resp.Queue, rideQueueResp{
							Name: ride.Name,
						})
					}
				}

				return c.JSON(200, resp)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
