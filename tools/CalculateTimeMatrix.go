package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type respBody struct {
	Durations [][]float64 `Json:"durations"`
}

var httpClient = &http.Client{Timeout: 2 * time.Second}

func sliceIndexMin(slice []float64) int {
	var min float64
	var index int

	if len(slice) > 0 {
		min = slice[0]
		index = 0
	}

	for i := 1; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
			index = i
		}
	}
	return index
}

func CalculateTimeMatrix(rides []models.Record, drivers []models.Record) int {
	base_url := "http://localhost:5000/table/v1/driving/"

	for _, ride := range rides {
		base_url = fmt.Sprintf("%s%s,%s;", base_url, ride.GetDataValue("origin_longitude"), ride.GetDataValue("origin_latitude"))
	}

	for i, driver := range drivers {
		if i != len(drivers)-1 {
			base_url = fmt.Sprintf("%s%s,%s;", base_url, driver.GetDataValue("current_longitude"), driver.GetDataValue("current_latitude"))
		} else {
			base_url = fmt.Sprintf("%s%s,%s", base_url, driver.GetDataValue("current_longitude"), driver.GetDataValue("current_latitude"))
		}
	}

	base_url = fmt.Sprintf("%s?sources=", base_url)

	for i := range drivers {
		if i != len(drivers)-1 {
			base_url = fmt.Sprintf("%s%d;", base_url, len(rides)+i)
		} else {
			base_url = fmt.Sprintf("%s%d", base_url, len(rides)+i)
		}
	}

	base_url = fmt.Sprintf("%s&destinations=", base_url)

	for i := range rides {
		if i != len(rides)-1 {
			base_url = fmt.Sprintf("%s%d;", base_url, i)
		} else {
			base_url = fmt.Sprintf("%s%d", base_url, i)
		}
	}

	var body respBody
	resp, err := httpClient.Get(base_url)
	if err != nil {
		fmt.Println("Error getting routes")
	}
	json.NewDecoder(resp.Body).Decode(&body)

	return sliceIndexMin(body.Durations[0])
}
