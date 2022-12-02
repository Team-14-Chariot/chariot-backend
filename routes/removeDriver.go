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

type removeDriverBody struct {
	DriverID string `json:"driver_id"`
}

func removeDriver(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/removeDriver",
		Handler: func(c echo.Context) error {
			var body removeDriverBody
			c.Bind(&body)

			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")

			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", body.DriverID)

			if driver != nil {
				if !driver.GetBoolDataValue("in_ride") {
					eventId := driver.GetStringDataValue("event_id")
					app.Dao().DeleteRecord(driver)

					helpers.UpdateDriverQueues(app, eventId, queues)

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
