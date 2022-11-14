package helpers

import (
	. "github.com/Team-14-Chariot/chariot-backend/models"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func GetNeededRides(app *pocketbase.PocketBase, rides_col *models.Collection, event_id string) []models.Record {
	rides := GetAllRecords(app, rides_col)

	n := 0
	for _, ride := range rides {
		if (ride.GetDataValue("event_id") == event_id) && (ride.GetDataValue("needs_ride") == true) {
			rides[n] = ride
			n++
		}
	}

	rides = rides[:n]
	return rides
}

func ConvertToRideObject(rides []models.Record) []*Ride {
	var ridesObjects []*Ride

	for _, ride := range rides {
		ridesObjects = append(ridesObjects, &Ride{
			ID:         ride.Id,
			Name:       ride.GetStringDataValue("rider_name"),
			GroupSize:  ride.GetIntDataValue("group_size"),
			OriginLat:  ride.GetStringDataValue("origin_latitude"),
			OriginLong: ride.GetStringDataValue("origin_longitude"),
			DestLat:    ride.GetStringDataValue("dest_latitude"),
			DestLong:   ride.GetStringDataValue("dest_longitude"),
			RideLength: ride.GetFloatDataValue("ride_length"),
			InRide:     ride.GetBoolDataValue("in_ride")})
	}

	return ridesObjects
}
