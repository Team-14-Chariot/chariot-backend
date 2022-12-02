package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type getDriverInfoRiderBody struct {
	RideID string `json:"ride_id"`
}

type driverInfoRiderResp struct {
	Name           string `json:"name"`
	CarDescription string `json:"car_description"`
	CarPlate       string `json:"car_license_plate"`
}

func getDriverInfoRider(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/getDriverInfoRider",
		Handler: func(c echo.Context) error {
			var body getDriverInfoRiderBody
			c.Bind(&body)

			rides, _ := app.Dao().FindCollectionByNameOrId("rides")
			ride, _ := app.Dao().FindFirstRecordByData(rides, "id", body.RideID)

			if ride != nil {
				drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
				driver, _ := app.Dao().FindFirstRecordByData(drivers, "id", ride.GetStringDataValue("driver_id"))

				if driver != nil {
					resp := driverInfoRiderResp{
						Name:           driver.GetStringDataValue("name"),
						CarDescription: driver.GetStringDataValue("car_description"),
						CarPlate:       driver.GetStringDataValue("car_licence_plate"),
					}

					return c.JSON(200, resp)
				}
			}
			return c.NoContent(400)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.RequireGuestOnly(),
		},
	})

	return nil
}
