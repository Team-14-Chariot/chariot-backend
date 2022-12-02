package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getDriverEventInfoBody struct {
	DriverID string `json:"driver_id"`
}

type driverEventInfoResp struct {
	RideMaxRadius int    `json:"maxRadius"`
	Address       string `json:"address"`
	EventName     string `json:"eventName"`
}

func getDriverEventInfo(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getDriverEventInfo",
		Handler: func(c echo.Context) error {
			var body getDriverEventInfoBody
			c.Bind(&body)

			drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", body.DriverID)

			if driver != nil {
				events, _ := app.Dao().FindCollectionByNameOrId("events")
				event, _ := app.Dao().FindFirstRecordByData(events, "event_id", driver.GetStringDataValue("event_id"))

				if event != nil {

					resp := driverEventInfoResp{
						RideMaxRadius: event.GetIntDataValue("ride_max_radius"),
						Address:       event.GetStringDataValue("address"),
						EventName:     event.GetStringDataValue("event_name"),
					}

					return c.JSON(200, resp)
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
