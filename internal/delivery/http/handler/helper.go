package handler

import (
	"net/http"

	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/validator"
)

// reusable validation helper
func ValidateRequest(w http.ResponseWriter, req interface{}) bool {
	if err := validator.ValidateStruct(req); err != nil {
		errMap := validator.ParseValidationError(err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("validation error", errMap))
		return false
	}
	return true
}
