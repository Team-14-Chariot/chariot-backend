package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type dropoffEarlyBody struct {
	RideID string `json:"ride_id"`
}

func dropoffEarly(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/dropoffEarly",
		Handler: func(c echo.Context) error {
			var body dropoffEarlyBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)
			ride.SetDataValue("in_ride", false)

			// Update queue and driver information here

			app.Dao().SaveRecord(ride)
			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
