package room

import (
	"encoding/json"
	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseRoom "hotel-loyalty/internal/usecase/room"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RoomHandler struct {
	roomUsecase *usecaseRoom.RoomUsecase
}

func NewRoomHandler(roomUsecase *usecaseRoom.RoomUsecase) *RoomHandler {
	return &RoomHandler{
		roomUsecase: roomUsecase,
	}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// 🔥 ambil hotel_id dari URL
	hotelIDStr := chi.URLParam(r, "hotel_id")

	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Room][Create] invalid hotel id:", hotelIDStr)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	var req request.CreateRoomRequest

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Println("[HANDLER][Room][Create] invalid request:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	// validation
	if !handler.ValidateRequest(w, req) {
		logger.ErrorLogger.Println("[HANDLER][Room][Create] validation failed")
		return
	}

	// mapping ke domain
	room := mapper.ToRoomDomain(req)
	room.HotelID = hotelID // 🔥 inject dari URL

	// call usecase
	resp, err := h.roomUsecase.Create(room)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Room][Create] error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println("[HANDLER][Room][Create] success:", req.Name)

	response.WriteJSON(w, http.StatusCreated, model.Success("Room created", resp))
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotel_id")

	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][Room][List] invalid hotel_id",
			"value", hotelIDStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	data, err := h.roomUsecase.List(hotelID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Room][List] failed",
			"hotelID", hotelID,
			"error", err,
		)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Room][List] success",
		"hotelID", hotelID,
		"count", len(data),
	)

	response.WriteJSON(w, http.StatusOK, model.Success("OK", data))
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][Room][Delete] invalid room_id",
			"value", idStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid room ID", nil))
		return
	}

	err = h.roomUsecase.Delete(roomID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Room][Delete] failed",
			"roomID", roomID,
			"error", err,
		)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Room][Delete] success",
		"roomID", roomID,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Room deleted", nil))
}

func (h *RoomHandler) SetInventory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	roomIDStr := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][Inventory] invalid room_id",
			"value", roomIDStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid room ID", nil))
		return
	}

	var req request.SetInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Inventory] invalid body",
			"error", err,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid start_date format (YYYY-MM-DD)", nil))
		return
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid end_date format (YYYY-MM-DD)", nil))
		return
	}

	err = h.roomUsecase.SetInventory(roomID, start, end, req.Stock)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Inventory] failed",
			"roomID", roomID,
			"error", err,
		)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Inventory] success",
		"roomID", roomID,
		"start", req.StartDate,
		"end", req.EndDate,
		"stock", req.Stock,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Inventory set", nil))
}
