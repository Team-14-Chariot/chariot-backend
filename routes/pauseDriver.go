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

type pauseDriveBody struct {
	DriverID string `json:"driver_id"`
}

func pauseDriver(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue, mutex *sync.RWMutex) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/pauseDriver",
		Handler: func(c echo.Context) error {
			var body pauseDriveBody
			c.Bind(&body)

			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")

			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", body.DriverID)

			if driver != nil {
				driver.SetDataValue("active", false)

				app.Dao().SaveRecord(driver)

				mutex.Lock()
				helpers.UpdateDriverQueues(app, driver.GetStringDataValue("event_id"), queues, nil, nil)
				mutex.Unlock()

				return c.NoContent(200)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
