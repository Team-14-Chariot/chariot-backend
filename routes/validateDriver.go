package routes

import (
    "net/http"

    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/models"
    )

type joinEventBody struct {
  EventCode      string `json:"event_id"`
    Name           string `json:"name"`
    CarCapacity    int    `json:"car_capacity"`
    CarDescription string `json:"car_description"`
    CarPlate       string `json:"car_licence_plate"`
    driverPass    string `json:"driver_password"`
}

func joinEvent(e *core.ServeEvent, app *pocketbase.PocketBase) error {
  e.Router.AddRoute(echo.Route{
Method: http.MethodPost,
Path:   "/api/validateDriverPassowrd",
Handler: func(c echo.Context) error {
var body joinEventBody

c.Bind(&body)

drivers, _ := app.Dao().FindCollectionByNameOrId("drivers")
events, _ := app.Dao().FindCollectionByNameOrId("events")
event, _ := app.Dao().FindFirstRecordByData(events, "event_id", body.EventCode)


driverPass, _ := app.Dao().FindFirstRecordByData(events, "driver_password", body.driverPass)

if driver != nil {
{if driverPass == body.driverPass} {
return true;
}


}
return false

},
Middlewares: []echo.MiddlewareFunc{
               apis.RequireGuestOnly(),
             },
             })

return nil
}

