package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func test_organizers(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/test_organizers",
		Handler: func(c echo.Context) error {
			collection, _ := app.Dao().FindCollectionByNameOrId("organizers")
			record, _ := app.Dao().FindFirstRecordByData(collection, "email", "hello@gmail.com")
			return c.String(200, fmt.Sprintf("Hello world: %v\n", record.GetDataValue("email")))
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
