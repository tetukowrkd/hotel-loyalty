package handler

import (
	"fmt"
	"net/http"

	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/validator"
)

func getType(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// reusable validation helper
func ValidateRequest(w http.ResponseWriter, req interface{}) bool {
	if err := validator.ValidateStruct(req); err != nil {
		logger.ErrorLogger.Println("[HANDLER][VALIDATION] validation failed:",
			"type=", getType(req),
			"error=", err)
		errMap := validator.ParseValidationError(err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("validation error", errMap))
		return false
	}
	return true
}
