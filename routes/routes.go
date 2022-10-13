package routes

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Routes(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return requestRide(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return joinEvent(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return leaveEvent(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return pauseDriver(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return resumeDriver(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return updateDriverStatus(e, app)
	})

	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return (addEventCode(e))
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return endEvent(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getRide(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return endRide(e, app)
	})
}
