package handler

import (
	"encoding/json"
	"net/http"

	"e-commerce/internal/domain"
	"e-commerce/internal/infrastructure/logger"
	"e-commerce/internal/model"
	"e-commerce/internal/pkg/errors"
	"e-commerce/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	// decode request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Register] Invalid Request:", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Error("Invalid Request Body", nil))
		return
	}

	// call usecase
	err = h.userUsecase.Register(&user)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Register] Error:", err)

		appErr, ok := err.(*errors.AppError)
		if ok {
			w.WriteHeader(appErr.Code)
			json.NewEncoder(w).Encode(model.Error(appErr.Message, nil))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Error("Internal Server Error", nil))
		return
	}

	// success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.Success("User Registered", nil))
}
