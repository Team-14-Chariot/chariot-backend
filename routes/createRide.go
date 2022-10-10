package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type createRideBody struct {
	rider            string
	driver           string
	origin_latitude  string
	origin_longitude string
	dest_latitude    string
	dest_longitude   string
	group_size       int
}

func createRide(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/createRide",
		Handler: func(c echo.Context) error {
			var body createRideBody
			c.Bind(&body)

			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
