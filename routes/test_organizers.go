package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func test_organizers(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "api/test_organizers",
		Handler: func(c echo.Context) error {
			// verify there is an organizer table
			collection, _ := app.Dao().FindCollectionByNameOrId("organizers")

			// get the first record from the table (email: hello@gmail.com | password: password)
			record, _ := app.Dao().FindFirstRecordByData(collection, "email", "hello@gmail.com")

			// return the organizer's information
			return c.String(200, "Hello organizer, "+record.GetStringDataValue("email")+". Your password is "+record.GetStringDataValue("password"))
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
