package routes

import (
	"net/http"
	"sync"

	"github.com/Team-14-Chariot/chariot-backend/helpers"
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type testBody struct {
	EventID string `json:"event_id"`
}

func test(e *core.ServeEvent, app *pocketbase.PocketBase, queues map[string]*DriverQueue, mutex *sync.RWMutex) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/test",
		Handler: func(c echo.Context) error {
			var body testBody
			c.Bind(&body)

			mutex.Lock()
			helpers.UpdateDriverQueues(app, body.EventID, queues, nil, nil)
			mutex.Unlock()

			return c.NoContent(200)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
