package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getEtaBody struct {
	RideID string `json:"ride_id"`
}

func getEta(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getEta",
		Handler: func(c echo.Context) error {
			var body getEtaBody
			c.Bind(&body)

		}
})

	

	return nil
}
