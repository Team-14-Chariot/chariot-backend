package models

type Ride struct {
	ID         string
	Name       string `json:"name"`
	GroupSize  int    `json:"groupSize"`
	OriginLat  string
	OriginLong string
	DestLat    string
	DestLong   string
	RideLength int
	InRide     bool `json:"inRide"`
}
