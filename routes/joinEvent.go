package routes

import (
	"net/http"

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
	CarPlate       string `json:"car_licence_plate"`
}

func joinEvent(e *core.ServeEvent, app *pocketbase.PocketBase) error {
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
				newDriver := models.NewRecord(drivers)

				newDriver.SetDataValue("name", body.Name)
				newDriver.SetDataValue("car_capacity", body.CarCapacity)
				newDriver.SetDataValue("car_description", body.CarDescription)
				newDriver.SetDataValue("car_licence_plate", body.CarPlate)
				newDriver.SetDataValue("event_id", body.EventCode)
				newDriver.SetDataValue("active", true)

				app.Dao().SaveRecord(newDriver)

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
