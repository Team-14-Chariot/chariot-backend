package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getEtaBody struct {
	RideID string `json:"ride_id"`
}

func getEta(e *core.ServeEvent, app core.App) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getEta",
		Handler: func(c echo.Context) error {
			var body getEtaBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)

			if ride != nil {
				return c.JSON(200, map[string]interface{}{"eta": ride.GetDataValue("eta")})
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
