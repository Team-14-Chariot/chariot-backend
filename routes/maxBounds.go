package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getMaxBounds struct {
	EventID string `json:"event_id"`
	
}

func getEta(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getMaxBounds",
		
	})

	Handler: func(c echo.Context) error {
		var body getMaxBounds
		c.Bind(&body)

		events, _ := app.Dao().FindCollectionByNameOrId("events")
		event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)
	

		if event != nil {
			return c.JSON(200, map[string]interface{}{
				"address": event.GetDataValue("address")
				"maxBounds": event.GetDataValue("ride_max_radius")
			})
		}

		return c.NoContent(400)
	},
	Middlewares: []echo.MiddlewareFunc{
		apis.RequireGuestOnly(),
	},

	return nil
}