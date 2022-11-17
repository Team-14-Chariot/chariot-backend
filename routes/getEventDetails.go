package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getEventDetailsBody struct {
	EventID string `json:"event_id"`
}

type eventDetailsResp struct {
	RideMaxRadius  int    `json:"maxRadius"`
	AcceptingRides bool   `json:"acceptingRides"`
	Address        string `json:"address"`
	EventName      string `json:"eventName"`
	RiderPassword  string `json:"ridePassword"`
	DriverPassword string `json:"driverPassword"`
}

func getEventdetails(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getEventDetails",
		Handler: func(c echo.Context) error {
			var body getEventDetailsBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)
			if event != nil {
				resp := eventDetailsResp{
					RideMaxRadius:  event.GetIntDataValue("ride_max_radius"),
					AcceptingRides: event.GetBoolDataValue("accept_rides"),
					Address:        event.GetStringDataValue("address"),
					EventName:      event.GetStringDataValue("event_name"),
					RiderPassword:  event.GetStringDataValue("rider_password"),
					DriverPassword: event.GetStringDataValue("driver_password"),
				}

				return c.JSON(200, resp)
			} else {
				return c.NoContent(400)
			}
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
