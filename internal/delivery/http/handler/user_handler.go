package handler

import (
	"encoding/json"
	"net/http"

	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/usecase"
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
	defer r.Body.Close()

	var req request.RegisterUserRequest

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Println("[HANDLER][Register] Invalid Request:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	// 🔥 pakai helper validation
	if !ValidateRequest(w, req) {
		return
	}

	// mapping ke domain
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// call usecase
	if err := h.userUsecase.Register(user); err != nil {
		logger.ErrorLogger.Println("[HANDLER][Register] Error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	// success
	response.WriteJSON(w, http.StatusCreated, model.Success("User Registered", nil))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request.LoginRequest

	// decode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	// validation
	if !ValidateRequest(w, req) {
		return
	}

	// usecase
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	// success
	response.WriteJSON(w, http.StatusOK, model.Success("Login Success", map[string]string{
		"token": token,
	}))
}
