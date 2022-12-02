package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type validateDriverPasswordBody struct {
	EventCode      string `json:"event_id"`
	DriverPassword string `json:"driver_password"`
}

func validateDriverPassword(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/validateDriverPassword",
		Handler: func(c echo.Context) error {
			var body validateDriverPasswordBody

			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventCode)

			if event != nil {
				if len(event.GetStringDataValue("driver_password")) > 0 {
					if event.GetStringDataValue("driver_password") != body.DriverPassword {
						return c.NoContent(400)
					}
				}

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
