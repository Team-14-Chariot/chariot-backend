package models

type Ride struct {
	ID         string
	Name       string `json:"name"`
	GroupSize  int    `json:"groupSize"`
	OriginLat  string
	OriginLong string
	DestLat    string
	DestLong   string
	RideLength float64
	InRide     bool `json:"inRide"`
	Edges      []Edge
}
