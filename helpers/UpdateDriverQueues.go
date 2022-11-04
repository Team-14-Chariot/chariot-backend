package helpers

import (
	. "github.com/Team-14-Chariot/chariot-backend/models"
	. "github.com/Team-14-Chariot/chariot-backend/tools"

	"github.com/pocketbase/pocketbase"
)

func UpdateDriverQueues(app *pocketbase.PocketBase, eventID string, queues map[string]*DriverQueue) {
	drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
	driversRecords := GetEventDrivers(app, drivers_col, eventID)
	rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
	ridesRecords := GetNeededRides(app, rides_col, eventID)
	drivers := ConvertToDriverObject(app, driversRecords)
	rides := ConvertToRideObject(ridesRecords)

	for i, driver := range drivers {
		queues[driver.ID] = InitDriverQueue()
		for _, ride := range rides {
			length := CalculateLength(driver.CurrentLat, driver.CurrentLong, ride.OriginLat, ride.OriginLong)
			drivers[i].Edges = append(driver.Edges, Edge{ID: ride.ID, Weight: length})
		}
	}

	for i, _ := range rides {
		for j, _ := range rides {
			if j > i {
				currRideToOtherLength := CalculateLength(rides[i].DestLat, rides[i].DestLong, rides[j].OriginLat, rides[j].OriginLong)
				rides[i].Edges = append(rides[i].Edges, Edge{ID: rides[j].ID, Weight: currRideToOtherLength})
				OtherRideToCurrLength := CalculateLength(rides[j].DestLat, rides[j].DestLong, rides[i].OriginLat, rides[j].OriginLong)
				rides[j].Edges = append(rides[j].Edges, Edge{ID: rides[i].ID, Weight: OtherRideToCurrLength})
			}
		}
	}

	assigned := []string{}
	i := 0
	for {
		if len(assigned) == len(rides) {
			break
		}

		startingNode := queues[drivers[i].ID]
		if startingNode.GetLastRide() == nil { // Driver has no rides yet
			edgeIndex, newEdge := getLowestDriverEdge(drivers[i])
			if edgeIndex != -1 {
				newRideIndex, newRide := findRide(rides, newEdge.ID)
				if isAssigned(assigned, newEdge.ID) {
					oldDriverQueue := queues[newRide.DriverID]

					oldTripTime := calculateTotalTripLength(oldDriverQueue, findDriver(drivers, newRide.DriverID), newEdge.ID)
					newTripTime := calculateTotalTripLength(startingNode, drivers[i], newEdge.ID) + newEdge.Weight

					if newTripTime < oldTripTime {
						oldDriverQueue.RemoveRide(&newRide)
						rides[newRideIndex].DriverID = drivers[i].ID
						startingNode.InsertRide(&newRide)
					} else {
						drivers[i].Edges = deleteDriverEdge(edgeIndex, drivers[i])
					}
				} else {
					rides[newRideIndex].DriverID = drivers[i].ID
					startingNode.InsertRide(&newRide)
					assigned = append(assigned, newRide.ID)
				}
			}
		} else {
			edgeIndex, newEdge := getLowestRideEdge(*startingNode.GetLastRide())
			if edgeIndex != -1 {
				newRideIndex, newRide := findRide(rides, newEdge.ID)
				if isAssigned(assigned, newEdge.ID) {
					oldDriverQueue := queues[newRide.DriverID]

					oldTripTime := calculateTotalTripLength(oldDriverQueue, findDriver(drivers, newRide.DriverID), newEdge.ID)
					newTripTime := calculateTotalTripLength(startingNode, drivers[i], newEdge.ID) + newEdge.Weight

					if newTripTime < oldTripTime {
						oldDriverQueue.RemoveRide(&newRide)
						rides[newRideIndex].DriverID = drivers[i].ID
						startingNode.InsertRide(&newRide)
					} else {
						rideIndex := findRideIndex(rides, *startingNode.GetLastRide())
						rides[rideIndex].Edges = deleteRideEdge(edgeIndex, rides[rideIndex])
					}
				} else {
					rides[newRideIndex].DriverID = drivers[i].ID
					startingNode.InsertRide(&newRide)
					assigned = append(assigned, newRide.ID)
				}
			}
		}

		i++
		if i == len(drivers) {
			i = 0
		}
	}
}

func getLowestDriverEdge(driver Driver) (int, Edge) {
	if len(driver.Edges) > 0 {
		lowestEdge := driver.Edges[0]
		lowestIndex := 0
		for i, edge := range driver.Edges {
			if edge.Weight < lowestEdge.Weight {
				lowestEdge = edge
				lowestIndex = i
			}
		}

		return lowestIndex, lowestEdge
	}
	return -1, Edge{}
}

func getLowestRideEdge(ride Ride) (int, Edge) {
	if len(ride.Edges) > 0 {
		lowestEdge := ride.Edges[0]
		lowestIndex := 0
		for i, edge := range ride.Edges {
			if edge.Weight < lowestEdge.Weight {
				lowestEdge = edge
				lowestIndex = i
			}
		}

		return lowestIndex, lowestEdge
	}
	return -1, Edge{}
}

func isAssigned(assigned []string, ID string) bool {
	for _, ride := range assigned {
		if ride == ID {
			return true
		}
	}

	return false
}

func findRide(rides []Ride, ID string) (int, Ride) {
	for i, ride := range rides {
		if ride.ID == ID {
			return i, ride
		}
	}

	return 0, Ride{}
}

func findDriver(drivers []Driver, ID string) Driver {
	for _, driver := range drivers {
		if driver.ID == ID {
			return driver
		}
	}

	return Driver{}
}

func calculateTotalTripLength(queue *DriverQueue, driver Driver, endingID string) float64 {
	totalLength := 0.0
	rides := queue.GetRides()

	i := 0
	for {
		if i > len(rides)-1 {
			break
		}

		if rides[i].ID != endingID {
			if i == 0 {
				totalLength += findDriverToFirstRideWeight(rides[0].ID, driver)
			} else if i < len(rides)-1 {
				totalLength += findRidetoRideEdgeWeight(rides[i+1].ID, rides[i])
			}
		}

		totalLength += rides[i].RideLength
		i++
	}

	return totalLength
}

func findDriverToFirstRideWeight(ID string, driver Driver) float64 {
	for _, edge := range driver.Edges {
		if edge.ID == ID {
			return edge.Weight
		}
	}
	return 0.0
}

func findRidetoRideEdgeWeight(ID string, ride Ride) float64 {
	for _, edge := range ride.Edges {
		if edge.ID == ID {
			return edge.Weight
		}
	}
	return 0.0
}

func deleteDriverEdge(index int, driver Driver) []Edge {
	driver.Edges[index] = driver.Edges[len(driver.Edges)-1]
	return driver.Edges[:len(driver.Edges)-1]
}

func deleteRideEdge(index int, ride Ride) []Edge {
	ride.Edges[index] = ride.Edges[len(ride.Edges)-1]
	return ride.Edges[:len(ride.Edges)-1]
}

func findRideIndex(rides []Ride, toFind Ride) int {
	for i, ride := range rides {
		if ride.ID == toFind.ID {
			return i
		}
	}
	return -1
}
