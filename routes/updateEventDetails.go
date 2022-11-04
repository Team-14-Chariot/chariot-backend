package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type updateEventDetailsBody struct {
	EventID       string `json:"event_id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	MaxRadius     int    `json:"max_radius"`
	RiderPassword string `json:"rider_password"`
}

func updateEventDetails(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/updateEventDetails",
		Handler: func(c echo.Context) error {
			var body updateEventDetailsBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)
			if event != nil {
				if len(body.Name) > 0 {
					event.SetDataValue("event_name", body.Name)
				}

				if len(body.Address) > 0 {
					event.SetDataValue("address", body.Address)
				}

				if len(body.RiderPassword) > 0 {
					event.SetDataValue("rider_password", body.RiderPassword)
				}

				if body.MaxRadius > 0 {
					event.SetDataValue("ride_max_radius", body.MaxRadius)
				}

				app.Dao().SaveRecord(event)

				return c.NoContent(200)
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
