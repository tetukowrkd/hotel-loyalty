package handler

import (
	"encoding/json"
	"net/http"

	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/infrastructure/middleware"
	"hotel-loyalty/internal/mapper"
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
	user := mapper.ToUserDomain(req)

	// call usecase
	resp, err := h.userUsecase.Register(user)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Register] Error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	// success
	logger.InfoLogger.Println("[HANDLER][Register] success:", req.Email)
	response.WriteJSON(w, http.StatusCreated, model.Success("User Registered", resp))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	if !ValidateRequest(w, req) {
		return
	}

	resp, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Login] error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}
	logger.InfoLogger.Println("[HANDLER][Login] success:", req.Email)
	response.WriteJSON(w, http.StatusOK, model.Success("Login Success", resp))
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if token == "" {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Missing token", nil))
		return
	}

	if err := h.userUsecase.Logout(token); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Logout failed", nil))
		return
	}

	response.WriteJSON(w, http.StatusOK, model.Success("Logout success", nil))
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	if !ValidateRequest(w, req) {
		return
	}

	resp, err := h.userUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	response.WriteJSON(w, http.StatusOK, model.Success("Token refreshed", resp))
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.Claims)
	if !ok {
		response.WriteJSON(w, http.StatusUnauthorized, model.Error("Unauthorized", nil))
		return
	}

	userID := claims.UserID

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	if user == nil {
		response.WriteJSON(w, http.StatusNotFound, model.Error("User Not Found", nil))
		return
	}

	// 🔥 mapping ke response DTO
	resp := mapper.ToUserResponse(user)

	response.WriteJSON(w, http.StatusOK, model.Success("Profile Fetched", resp))
}
