package hotel

import (
	"encoding/json"
	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseHotel "hotel-loyalty/internal/usecase/hotel"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HotelHandler struct {
	hotelUsecase *usecaseHotel.HotelUsecase
}

func NewHotelHandler(hotelUsecase *usecaseHotel.HotelUsecase) *HotelHandler {
	return &HotelHandler{
		hotelUsecase: hotelUsecase,
	}
}

func (h *HotelHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request.CreateHotelRequest

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Println("[HANDLER][Hotel][Create] invalid request:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid request body", nil))
		return
	}

	// validation
	if !handler.ValidateRequest(w, req) {
		logger.ErrorLogger.Println("[HANDLER][Hotel][Create] validation failed")
		return
	}

	// mapping ke domain (CONSISTENT sama user)
	hotel := mapper.ToHotelDomain(req)

	// call usecase
	resp, err := h.hotelUsecase.Create(hotel)
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][Hotel][Create] error:", err)

		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	// success
	logger.InfoLogger.Println("[HANDLER][Hotel][Create] success:", req.Name)

	response.WriteJSON(w, http.StatusCreated, model.Success("Hotel created", resp))
}

func (h *HotelHandler) List(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")

	starStr := r.URL.Query().Get("star")
	star := 0
	if starStr != "" {
		var err error
		star, err = strconv.Atoi(starStr)
		if err != nil {
			logger.InfoLogger.Println(
				"[HANDLER][Hotel][List] invalid star param",
				"star", starStr,
			)

			response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid star parameter", nil))
			return
		}
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, total, err := h.hotelUsecase.List(city, star, limit, offset)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Hotel][List] failed",
			"city", city,
			"star", star,
			"error", err,
		)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Hotel][List] success",
		"city", city,
		"star", star,
		"count", len(data),
		"total", total,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("OK", map[string]interface{}{
		"items": data,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}))
}

func (h *HotelHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "hotel_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][Hotel][GetByID] invalid id",
			"id", idStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid ID", nil))
		return
	}

	data, err := h.hotelUsecase.GetByID(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Hotel][GetByID] failed",
			"id", id,
			"error", err,
		)

		response.WriteJSON(w, http.StatusNotFound, model.Error("Not Found", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Hotel][GetByID] success",
		"id", id,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("OK", data))
}

func (h *HotelHandler) Publish(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "hotel_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.InfoLogger.Println(
			"[HANDLER][Hotel][Publish] invalid id",
			"id", idStr,
		)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid ID", nil))
		return
	}

	err = h.hotelUsecase.Publish(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][Hotel][Publish] failed",
			"id", id,
			"error", err,
		)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][Hotel][Publish] success",
		"id", id,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Hotel published", nil))
}

func (h *HotelHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "hotel_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid ID", nil))
		return
	}

	err = h.hotelUsecase.Delete(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.WriteJSON(w, appErr.Code, model.Error(appErr.Message, nil))
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("Internal Server Error", nil))
		return
	}

	response.WriteJSON(w, http.StatusOK, model.Success("Hotel deleted", nil))
}
