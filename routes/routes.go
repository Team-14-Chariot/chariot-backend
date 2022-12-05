package routes

import (
	"sync"

	. "github.com/Team-14-Chariot/chariot-backend/models"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Routes(app *pocketbase.PocketBase, queues map[string]*DriverQueue) {
	mutex := &sync.RWMutex{}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return requestRide(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return joinEvent(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return leaveEvent(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return pauseDriver(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return resumeDriver(e, app, queues, mutex)
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
		return getRide(e, app, queues)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return endRide(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getEta(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return validateEvent(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return dropoffEarly(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return pickupEarly(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getRideQueues(e, app, queues)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getEventdetails(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return updateEventDetails(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getEventDrivers(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return removeDriver(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getRideQueue(e, app, queues)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getDriverEventInfo(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return validateDriverPassword(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getDriverInfoRider(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return cancelRide(e, app, queues, mutex)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getBallparkETA(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return getRouteETA(e, app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return test(e, app, queues, mutex)
	})
}
