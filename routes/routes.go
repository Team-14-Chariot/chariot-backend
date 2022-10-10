package routes

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Routes(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return test_organizers(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return test(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return createRide(e, app)
	})

	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return (addEventCode(e))
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return endEvent(e, app)
	})

}
