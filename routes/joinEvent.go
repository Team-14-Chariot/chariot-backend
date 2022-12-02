package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type joinEventBody struct {
	EventCode      string `json:"event_id"`
	Name           string `json:"name"`
	CarCapacity    int    `json:"car_capacity"`
	CarDescription string `json:"car_description"`
	CarPlate       string `json:"car_license_plate"`
	DriverPassword string `json:"driver_password"`
}

func joinEvent(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/joinEvent",
		Handler: func(c echo.Context) error {
			var body joinEventBody

			c.Bind(&body)

			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventCode)
			if event != nil {
				if len(event.GetStringDataValue("driver_password")) > 0 {
					if event.GetStringDataValue("driver_password") != body.DriverPassword {
						return c.NoContent(400)
					}
				}

				newDriver := models.NewRecord(drivers)

				newDriver.SetDataValue("name", body.Name)
				newDriver.SetDataValue("car_capacity", body.CarCapacity)
				newDriver.SetDataValue("car_description", body.CarDescription)
				newDriver.SetDataValue("car_licence_plate", body.CarPlate)
				newDriver.SetDataValue("event_id", body.EventCode)
				newDriver.SetDataValue("active", true)
				newDriver.SetDataValue("in_ride", false)
				newDriver.SetDataValue("has_rider", false)

				app.Dao().SaveRecord(newDriver)

				helpers.UpdateDriverQueues(app, body.EventCode, queues, nil)

				return c.JSON(200, map[string]interface{}{"driver_id": newDriver.Id})
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
