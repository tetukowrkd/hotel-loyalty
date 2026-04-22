package hotel

import (
	"encoding/json"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseHotel "hotel-loyalty/internal/usecase/hotel"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HotelFacilityHandler struct {
	hotelFacilityUsecase *usecaseHotel.HotelFacilityUsecase
}

func NewHotelFacilityHandler(hotelFacilityUsecase *usecaseHotel.HotelFacilityUsecase) *HotelFacilityHandler {
	return &HotelFacilityHandler{
		hotelFacilityUsecase: hotelFacilityUsecase,
	}
}

func (h *HotelFacilityHandler) Assign(w http.ResponseWriter, r *http.Request) {

	hotelIDStr := chi.URLParam(r, "hotel_id")
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][HotelFacility][Assign] invalid hotel id",
			"hotelID", hotelIDStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	var req request.AssignFacilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][HotelFacility][Assign] invalid request body",
			"hotelID", hotelID,
			"error", err,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid body", nil))
		return
	}

	err = h.hotelFacilityUsecase.Assign(hotelID, req.FacilityIDs)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelFacility][Assign] assign failed",
			"hotelID", hotelID,
			"facilityCount", len(req.FacilityIDs),
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
		"[HANDLER][HotelFacility][Assign] success",
		"hotelID", hotelID,
		"facilityCount", len(req.FacilityIDs),
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Facilities assigned", nil))
}

func (h *HotelFacilityHandler) GetByHotel(w http.ResponseWriter, r *http.Request) {

	hotelIDStr := chi.URLParam(r, "hotel_id")
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][HotelFacility][GetByHotel] invalid hotel id",
			"hotelID", hotelIDStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	data, err := h.hotelFacilityUsecase.GetByHotelID(hotelID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelFacility][GetByHotel] failed",
			"hotelID", hotelID,
			"error", err,
		)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][HotelFacility][GetByHotel] success",
		"hotelID", hotelID,
		"count", len(data),
	)

	response.WriteJSON(w, http.StatusOK, model.Success("OK", data))
}

func (h *HotelFacilityHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	data, err := h.hotelFacilityUsecase.GetAll()
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelFacility][GetAll] failed",
			"error", err,
		)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	response.WriteJSON(w, http.StatusOK, model.Success("OK", data))
}

func (h *HotelFacilityHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req request.CreateFacilityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][HotelFacility][Create] invalid request body",
			"error", err,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid body", nil))
		return
	}

	id, err := h.hotelFacilityUsecase.Create(req.Name, req.Icon)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelFacility][Create] failed",
			"name", req.Name,
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
		"[HANDLER][HotelFacility][Create] success",
		"id", id,
		"name", req.Name,
	)

	response.WriteJSON(w, http.StatusCreated, model.Success(
		"Facility created",
		response.CreateFacilityResponse{
			ID:   id,
			Name: req.Name,
			Icon: req.Icon,
		},
	))
}

func (h *HotelFacilityHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "facility_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][HotelFacility][Delete] invalid facility id",
			"facilityID", idStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid facility ID", nil))
		return
	}

	err = h.hotelFacilityUsecase.Delete(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelFacility][Delete] failed",
			"facilityID", id,
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
		"[HANDLER][HotelFacility][Delete] success",
		"facilityID", id,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Facility deleted", nil))
}
