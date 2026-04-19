package identity

import (
	"encoding/json"
	"net/http"

	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseIdentity "hotel-loyalty/internal/usecase/identity"
)

type TokenHandler struct {
	tokenUsecase *usecaseIdentity.TokenUsecase
}

func NewTokenHandler(tokenUsecase *usecaseIdentity.TokenUsecase) *TokenHandler {
	return &TokenHandler{
		tokenUsecase: tokenUsecase,
	}
}

func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Println("[HANDLER][RefreshToken] invalid body:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid Request Body", nil))
		return
	}

	if !handler.ValidateRequest(w, req) {
		logger.ErrorLogger.Println("[HANDLER][RefreshToken] validation failed")
		return
	}

	resp, err := h.tokenUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][RefreshToken] error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println("[HANDLER][RefreshToken] success")

	response.WriteJSON(w, http.StatusOK, model.Success("Token Refreshed", resp))
}
