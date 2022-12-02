package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getBallparkETABody struct {
	EventID string `json:"event_id"`
}

func getBallparkETA(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getBallparkETA",
		Handler: func(c echo.Context) error {
			var body getBallparkETABody
			c.Bind(&body)

			rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
			driver_col, _ := app.Dao().FindCollectionByNameOrId("drivers")

			drivers := helpers.GetEventDrivers(app, driver_col, body.EventID)
			rides := helpers.GetNeededRides(app, rides_col, body.EventID)

			eta := len(rides) * 1200 / len(drivers)
			return c.JSON(200, map[string]interface{}{"eta": eta})
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
