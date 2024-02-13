package models

type ParkingDTO struct {
	Car string `json:"car,omitempty"`
}

type Parking struct {
	Id  int    `json:"id,omitempty"`
	Car string `json:"car,omitempty"`
}
