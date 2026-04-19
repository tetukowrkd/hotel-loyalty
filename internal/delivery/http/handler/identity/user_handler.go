package identity

import (
	"encoding/json"
	"net/http"

	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/infrastructure/middleware"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseIdentity "hotel-loyalty/internal/usecase/identity"
)

type UserHandler struct {
	userUsecase *usecaseIdentity.UserUsecase
}

func NewUserHandler(userUsecase *usecaseIdentity.UserUsecase) *UserHandler {
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
	if !handler.ValidateRequest(w, req) {
		logger.ErrorLogger.Println("[HANDLER][RefreshToken] validation failed")
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
		logger.ErrorLogger.Println("[HANDLER][RefreshToken] invalid body:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid Request Body", nil))
		return
	}

	if !handler.ValidateRequest(w, req) {
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
		logger.ErrorLogger.Println("[HANDLER][Logout] missing token")

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Missing Token", nil))
		return
	}

	if err := h.userUsecase.Logout(token); err != nil {
		logger.ErrorLogger.Println("[HANDLER][Logout] error:", err)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Logout Failed", nil))
		return
	}

	logger.InfoLogger.Println("[HANDLER][Logout] success")

	response.WriteJSON(w, http.StatusOK, model.Success("Logout Success", nil))
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.Claims)
	if !ok {
		logger.ErrorLogger.Println("[HANDLER][Profile] unauthorized access")

		response.WriteJSON(w, http.StatusUnauthorized, model.Error("Unauthorized", nil))
		return
	}

	userID := claims.UserID

	logger.InfoLogger.Println("[HANDLER][Profile] request:", userID)

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Profile] error:", err)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	if user == nil {
		logger.InfoLogger.Println("[HANDLER][Profile] user not found:", userID)

		response.WriteJSON(w, http.StatusNotFound, model.Error("User Not Found", nil))
		return
	}

	// 🔥 mapping ke response DTO
	resp := mapper.ToUserResponse(user)

	logger.InfoLogger.Println("[HANDLER][Profile] success:", userID)

	response.WriteJSON(w, http.StatusOK, model.Success("Profile Fetched", resp))
}

func (h *UserHandler) GetMemberQR(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("[HANDLER][GetMemberQR] hit")

	claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.Claims)
	if !ok {
		logger.ErrorLogger.Println("[HANDLER][GetMemberQR] unauthorized")
		response.WriteJSON(w, http.StatusUnauthorized, model.Error("Unauthorized", nil))
		return
	}

	userID := claims.UserID
	logger.InfoLogger.Println("[HANDLER][GetMemberQR] userID:", userID)

	qr, err := h.userUsecase.GetMemberQR(userID)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][GetMemberQR] error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	resp := response.QRResponse{
		QRCode: qr,
	}

	logger.InfoLogger.Println("[HANDLER][GetMemberQR] success userID:", userID)

	response.WriteJSON(w, http.StatusOK, model.Success("QR Generated", resp))
}
