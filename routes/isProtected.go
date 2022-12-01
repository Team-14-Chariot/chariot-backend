package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type isProtectedBody struct {
	EventCode      string `json:"event_id"`
	Name           string `json:"name"`
	CarCapacity    int    `json:"car_capacity"`
	CarDescription string `json:"car_description"`
	CarPlate       string `json:"car_licence_plate"`
	EventPass      string `json:"driver_password"`
}

func isProtected(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/isProtected",
		Handler: func(c echo.Context) error {
			var body isProtectedBody

			c.Bind(&body)

			// drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
			// events, _ := app.Dao().FindCollectionByNameOrId("events")
			// event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventCode)
			// password, _ := app.Dao().FindFirstRecordByData(events, "password", body.EventPass)
			// if event != nil {
			// 	if password != nil {
			// 		return true
			// 	}
			// 	return false

			// }
			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
