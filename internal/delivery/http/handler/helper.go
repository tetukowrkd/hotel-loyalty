package handler

import (
	"net/http"

	"e-commerce/internal/model"
	"e-commerce/internal/model/response"
	"e-commerce/internal/pkg/validator"
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
