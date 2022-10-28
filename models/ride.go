package models

type Ride struct {
	ID         string
	GroupSize  int
	OriginLat  string
	OriginLong string
	DestLat    string
	DestLong   string
	RideLength int
}
