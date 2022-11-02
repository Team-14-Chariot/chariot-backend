package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type EndEventBody struct {
	Event_id string `json:"event_id"`
}

func endEvent(e *core.ServeEvent, app core.App) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/endEvent",
		Handler: func(c echo.Context) error {
			var body EndEventBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			if events != nil {
				record, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.Event_id)

				if record != nil {
					record.SetDataValue("accept_rides", false)

					app.Dao().SaveRecord(record)
					return c.NoContent(200)
				} else {
					return c.NoContent(400)
				}
			} else {
				return c.NoContent(400)
			}
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
