package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type endEventBody struct {
	Event_id string `json:"event_id"`
}

func endEvent(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/endEvent",
		Handler: func(c echo.Context) error {
			var body endEventBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			record, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.Event_id)
			record.SetDataValue("accept_rides", false)

			app.Dao().SaveRecord(record)
			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
