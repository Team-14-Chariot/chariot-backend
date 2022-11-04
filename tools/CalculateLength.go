package tools

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

type lengthRespBody struct {
	Code  string      `json:"code"`
	Route []CalRoutes `json:"routes"`
}

type CalRoutes struct {
	Length float64 `json:"duration"`
}

func CalculateLength(startLat string, startLong string, endLat string, endLong string) float64 {
	if len(startLat) == 0 || len(startLong) == 0 || len(endLat) == 0 || len(endLong) == 0 {
		return 0
	} else {
		url := fmt.Sprintf("http://localhost:5000/route/v1/driving/%s,%s;%s,%s", startLong, startLat, endLong, endLat)

		var body lengthRespBody
		resp, err := httpClient.Get(url)
		if err != nil {
			fmt.Println("Error getting routes")
		}

		json.NewDecoder(resp.Body).Decode(&body)
		return body.Route[0].Length
	}
}

func CalculateRideLength(ride models.Record) float64 {
	return CalculateLength(ride.GetStringDataValue("origin_latitude"), ride.GetStringDataValue("origin_longitude"), ride.GetStringDataValue("dest_latitude"), ride.GetStringDataValue("dest_longitude"))
}
