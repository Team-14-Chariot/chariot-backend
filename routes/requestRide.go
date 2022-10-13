package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type requestRideBody struct {
	EventID    string `json:"event_id"`
	OriginLat  string `json:"origin_latitude"`
	OriginLong string `json:"origin_longitude"`
	DestLat    string `json:"dest_latitude"`
	DestLong   string `json:"dest_longitude"`
	GroupSize  int    `json:"group_size"`
}

func requestRide(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/requestRide",
		Handler: func(c echo.Context) error {
			var body requestRideBody
			c.Bind(&body)

			drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
			drivers := helpers.GetAllRecords(app, drivers_col)

			if len(drivers) > 0 {
				// 	rides, _ := app.Dao().FindCollectionByNameOrId("rides")

				// 	events, _ := app.Dao().FindCollectionByNameOrId("events")
				// 	event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)

				// 	if (event != nil) && (event.GetDataValue("accept_rides") == true) {
				// 		newRide := models.NewRecord(rides)

				// 		app.Dao().SaveRecord(newRide)

				// 		return c.JSON(200, map[string]interface{}{"ride_id": newRide.Id})
				// 	}
				return c.NoContent(200)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
