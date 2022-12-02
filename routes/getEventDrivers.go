package routes

import (
	"net/http"

	"github.com/Team-14-Chariot/chariot-backend/helpers"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getEventDriversBody struct {
	EventID string `json:"event_id"`
}

type getEventDriversResp struct {
	Drivers []driverResp `json:"drivers"`
}

type driverResp struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CarCapacity    int    `json:"car_capacity"`
	CarDescription string `json:"car_description"`
	CarPlate       string `json:"car_plate"`
	Active         bool   `json:"active"`
	CurrentLat     string `json:"current_latitude"`
	CurrentLong    string `json:"current_longitude"`
	InRide         bool   `json:"in_ride"`
	HasRider       bool   `json:"has_rider"`
	RidesCompleted int    `json:"rides_completed"`
}

func getEventDrivers(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getEventDrivers",
		Handler: func(c echo.Context) error {
			var body getEventDriversBody
			c.Bind(&body)

			events, _ := app.Dao().FindCollectionByNameOrId("events")
			event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventID)

			if event != nil {
				drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
				drivers := helpers.GetAllEventDrivers(app, drivers_col, body.EventID)

				resp := getEventDriversResp{Drivers: []driverResp{}}

				for _, driver := range drivers {
					resp.Drivers = append(resp.Drivers, driverResp{
						ID:             driver.Id,
						Name:           driver.GetStringDataValue("name"),
						CarCapacity:    driver.GetIntDataValue("car_capacity"),
						CarDescription: driver.GetStringDataValue("car_description"),
						CarPlate:       driver.GetStringDataValue("car_licence_plate"),
						Active:         driver.GetBoolDataValue("active"),
						CurrentLat:     driver.GetStringDataValue("current_latitude"),
						CurrentLong:    driver.GetStringDataValue("current_longitude"),
						InRide:         driver.GetBoolDataValue("in_ride"),
						HasRider:       driver.GetBoolDataValue("has_rider"),
						RidesCompleted: driver.GetIntDataValue("ride_count"),
					})
				}

				return c.JSON(200, resp)
			}

			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
