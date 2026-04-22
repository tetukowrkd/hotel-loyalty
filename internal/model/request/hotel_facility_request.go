package request

import "github.com/google/uuid"

type AssignFacilityRequest struct {
	FacilityIDs []uuid.UUID `json:"facility_ids"`
}

type CreateFacilityRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
