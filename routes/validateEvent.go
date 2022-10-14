package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type validateEventBody struct {
	EventID string `json:"event_id"`
}

func validateEvent(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/validateEvent",
		Handler: func(c echo.Context) error {
			var body validateEventBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)

			if event != nil {
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
