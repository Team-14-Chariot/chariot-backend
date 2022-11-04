package models

type Driver struct {
	ID          string
	Capacity    int
	CurrentLat  string
	CurrentLong string
	Edges       []Edge
}
