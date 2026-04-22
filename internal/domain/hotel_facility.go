package domain

import "github.com/google/uuid"

type HotelFacility struct {
	ID   uuid.UUID
	Name string
	Icon string
}

type HotelFacilityMap struct {
	HotelID    uuid.UUID
	FacilityID uuid.UUID
}
