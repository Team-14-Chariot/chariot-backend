package helpers

import (
	"encoding/json"
	"fmt"

	. "github.com/Team-14-Chariot/chariot-backend/models"
	. "github.com/Team-14-Chariot/chariot-backend/tools"

	"github.com/pocketbase/pocketbase"
)

func UpdateDriverQueues(app *pocketbase.PocketBase, eventID string, queues map[string]*DriverQueue) {
	drivers_col, _ := app.Dao().FindCollectionByNameOrId("drivers")
	driversRecords := GetEventDrivers(app, drivers_col, eventID)
	rides_col, _ := app.Dao().FindCollectionByNameOrId("rides")
	ridesRecords := GetNeededRides(app, rides_col, eventID)
	drivers := ConvertToDriverObject(app, driversRecords) // List of pointers to all active drivers for the given event
	rides := ConvertToRideObject(ridesRecords)            // List of pointers to all rides that have not been picked up for the given event
	ridesMap := make(map[string]*Ride)                    // Hashmap for finding ride pointers given a ride id to speed up finding rides from edges
	driversMap := make(map[string]*Driver)                // Hashmap for finding driver points given a driver id to speed up finding drivers when they have already been assigned a ride

	// Going through all drivers and initilizing their queues or resetting them if they already existed
	// Also adds edges to all rides with weights representing drive time from drivers location or where their ongoing ride will end
	for i, driver := range drivers {
		queues[driver.ID] = InitDriverQueue()
		for _, ride := range rides {
			length := CalculateLength(driver.CurrentLat, driver.CurrentLong, ride.OriginLat, ride.OriginLong)
			drivers[i].Edges = append(driver.Edges, Edge{ID: ride.ID, Weight: length})
		}

		driversMap[drivers[i].ID] = drivers[i]
	}

	// Going through all rides and adding edges with drive times between the start and ends of rides
	for i, _ := range rides {
		for j, _ := range rides {
			if j > i {
				currRideToOtherLength := CalculateLength(rides[i].DestLat, rides[i].DestLong, rides[j].OriginLat, rides[j].OriginLong)
				rides[i].Edges = append(rides[i].Edges, Edge{ID: rides[j].ID, Weight: currRideToOtherLength})
				OtherRideToCurrLength := CalculateLength(rides[j].DestLat, rides[j].DestLong, rides[i].OriginLat, rides[i].OriginLong)
				rides[j].Edges = append(rides[j].Edges, Edge{ID: rides[i].ID, Weight: OtherRideToCurrLength})
			}
		}

		ridesMap[rides[i].ID] = rides[i] // Adding the ride to the map for faster lookup
	}

	// for _, driver := range drivers {
	// 	out, _ := json.MarshalIndent(driver, "", "  ")
	// 	fmt.Println(string(out))
	// }

	// for _, ride := range rides {
	// 	out, _ := json.MarshalIndent(ride, "", "  ")
	// 	fmt.Println(string(out))
	// }

	assigned := []string{}
	i := 0 // Index of the driver who's queue is being added to
	anyChanges := false
	for {
		anyChanges = false

		driverQueue := queues[drivers[i].ID]
		if driverQueue.GetLastRide() == nil { // Driver has no rides in their queue
			edgeIndex, shortestEdge := findShortestDriverEdge(*drivers[i])

			if edgeIndex != -1 { // This should always be the case
				rideToAssign := ridesMap[shortestEdge.ID] // The ride the edge points to

				if isAssigned(assigned, rideToAssign.ID) { // If the ride is already assigned to another driver
					oldDriverQueue := queues[rideToAssign.DriverID]

					oldTripTime := calculateTotalTripLength(oldDriverQueue, *driversMap[rideToAssign.DriverID], rideToAssign.ID)
					newTripTime := calculateTotalTripLength(driverQueue, *drivers[i], rideToAssign.ID) + shortestEdge.Weight
					out1, _ := json.MarshalIndent(oldDriverQueue.GetRides(), "", "  ")
					fmt.Println(string(out1))
					out2, _ := json.MarshalIndent(oldDriverQueue.GetRides(), "", "  ")
					fmt.Println(string(out2))
					fmt.Printf("Index: %d, old trip time: %f, new trip time: %f\n", i, oldTripTime, newTripTime)

					if newTripTime < oldTripTime { // Remove the ride from the driver queue it was assigned to and add it to the current driver queue
						oldDriverQueue.RemoveRide(rideToAssign)
						rideToAssign.DriverID = drivers[i].ID
						driverQueue.InsertRide(rideToAssign)
						anyChanges = true
						fmt.Printf("Here5, index %d\n", i)
					} else {
						fmt.Printf("Here3, index %d\n", i)
						drivers[i].Edges = deleteDriverEdge(edgeIndex, *drivers[i])
						anyChanges = true
					}
				} else { // The ride has not already been assigned
					fmt.Printf("Here6, index %d\n", i)
					rideToAssign.DriverID = drivers[i].ID
					driverQueue.InsertRide(rideToAssign)
					assigned = append(assigned, rideToAssign.ID)
					anyChanges = true
				}
			} else {
				fmt.Printf("Here1, index %d\n", i)
			}
		} else { // The driver already has rides in their queue
			edgeIndex, shortestEdge := getShortestRideEdge(*driverQueue.GetLastRide())

			if edgeIndex != -1 {
				rideToAssign := ridesMap[shortestEdge.ID]

				if isAssigned(assigned, rideToAssign.ID) {
					oldDriverQueue := queues[rideToAssign.DriverID]

					oldTripTime := calculateTotalTripLength(oldDriverQueue, *driversMap[rideToAssign.DriverID], rideToAssign.ID)
					newTripTime := calculateTotalTripLength(driverQueue, *drivers[i], rideToAssign.ID) + shortestEdge.Weight
					out1, _ := json.MarshalIndent(oldDriverQueue.GetRides(), "", "  ")
					fmt.Println(string(out1))
					out2, _ := json.MarshalIndent(driverQueue.GetRides(), "", "  ")
					fmt.Println(string(out2))
					fmt.Printf("Index: %d, old trip time: %f, new trip time: %f\n", i, oldTripTime, newTripTime)

					if newTripTime < oldTripTime {
						oldDriverQueue.RemoveRide(rideToAssign)
						rideToAssign.DriverID = drivers[i].ID
						driverQueue.InsertRide(rideToAssign)
						anyChanges = true
						fmt.Printf("Here7, index %d\n", i)
					} else {
						fmt.Printf("Here4, index %d\n", i)
						driverQueue.GetLastRide().Edges = deleteRideEdge(edgeIndex, *driverQueue.GetLastRide())
						anyChanges = true
					}
				} else {
					fmt.Printf("Here8, index %d\n", i)
					rideToAssign.DriverID = drivers[i].ID
					driverQueue.InsertRide(rideToAssign)
					assigned = append(assigned, rideToAssign.ID)
					anyChanges = true
				}
			} else {
				fmt.Printf("Here2, index %d\n", i)
			}
		}

		i++
		if i == len(drivers) {
			if len(assigned) == len(rides) && !anyChanges {
				break
			}

			i = 0
		}
	}
	fmt.Println("----AFTER----")
	for _, driver := range drivers {
		out, _ := json.MarshalIndent(driver, "", "  ")
		fmt.Println(string(out))
	}

	for _, ride := range rides {
		out, _ := json.MarshalIndent(ride, "", "  ")
		fmt.Println(string(out))
	}
}

func findShortestDriverEdge(driver Driver) (int, Edge) {
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

func getShortestRideEdge(ride Ride) (int, Edge) {
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

func calculateTotalTripLength(queue *DriverQueue, driver Driver, endingID string) float64 {
	totalLength := 0.0
	rides := queue.GetRides()

	i := 0
	for {
		if i > len(rides)-1 {
			break
		}

		if i == 0 {
			totalLength += findDriverToFirstRideEdgeWeight(rides[i].ID, driver)
		} else {
			totalLength += findRidetoRideEdgeWeight(rides[i-1].ID, rides[i])
		}

		if rides[i].ID == endingID {
			break
		}

		totalLength += rides[i].RideLength

		i++
	}

	return totalLength
}

func findDriverToFirstRideEdgeWeight(rideID string, driver Driver) float64 {
	for _, edge := range driver.Edges {
		if edge.ID == rideID {
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
