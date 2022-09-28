package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func test(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/test",
		Handler: func(c echo.Context) error {
			collection, _ := app.Dao().FindCollectionByNameOrId("drivers")
			record, _ := app.Dao().FindFirstRecordByData(collection, "name", "Greg")
			return c.String(200, fmt.Sprintf("Hello world: %v\n", record.GetDataValue("car_licence_plate")))
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
