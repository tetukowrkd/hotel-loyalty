package hotel

import (
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	usecaseImageHotel "hotel-loyalty/internal/usecase/hotel"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HotelImageHandler struct {
	hotelImageUsecase *usecaseImageHotel.HotelImageUsecase
}

func NewHotelImageHandler(hotelImageUsecase *usecaseImageHotel.HotelImageUsecase) *HotelImageHandler {
	return &HotelImageHandler{
		hotelImageUsecase: hotelImageUsecase,
	}
}

func (h *HotelImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotel_id")

	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.InfoLogger.Println("[HANDLER][HotelImage] invalid hotel id:", hotelIDStr)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		logger.ErrorLogger.Println("[HANDLER][HotelImage] parse form error:", err)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid form data", nil))
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("No images provided", nil))
		return
	}

	var uploaded []*response.CreateHotelImageResponse

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			logger.ErrorLogger.Println(
				"[HANDLER][HotelImage] open file error",
				"hotelID", hotelID,
				"filename", header.Filename,
				"error", err,
			)
			continue
		}

		resp, err := h.hotelImageUsecase.Create(
			hotelID,
			file,
			header.Filename,
			header.Size,
		)

		file.Close()

		if err != nil {
			logger.ErrorLogger.Println(
				"[HANDLER][HotelImage] upload error",
				"hotelID", hotelID,
				"filename", header.Filename,
				"error", err,
			)
			continue
		}

		uploaded = append(uploaded, resp)
	}

	if len(uploaded) == 0 {
		logger.ErrorLogger.Println("[HANDLER][HotelImage] all uploads failed", "hotelID", hotelID)

		response.WriteJSON(w, http.StatusInternalServerError, model.Error("All uploads failed", nil))
		return
	}

	logger.InfoLogger.Println(
		"[HANDLER][HotelImage] upload success",
		"hotelID", hotelID,
		"count", len(uploaded),
	)

	response.WriteJSON(w, http.StatusCreated, model.Success("Images uploaded", uploaded))
}

func (h *HotelImageHandler) SetPrimary(w http.ResponseWriter, r *http.Request) {

	hotelIDStr := chi.URLParam(r, "hotel_id")
	imageIDStr := chi.URLParam(r, "image_id")

	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		logger.InfoLogger.Println("[HANDLER][HotelImage] invalid hotel id:", hotelIDStr)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		logger.InfoLogger.Println("[HANDLER][HotelImage] invalid image id:", imageIDStr)

		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid image ID", nil))
		return
	}

	err = h.hotelImageUsecase.SetPrimary(hotelID, imageID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelImage] set primary error",
			"hotelID", hotelID,
			"imageID", imageID,
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
		"[HANDLER][HotelImage] primary updated",
		"hotelID", hotelID,
		"imageID", imageID,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Primary image updated", nil))
}

func (h *HotelImageHandler) Delete(w http.ResponseWriter, r *http.Request) {

	hotelIDStr := chi.URLParam(r, "hotel_id")
	imageIDStr := chi.URLParam(r, "image_id")

	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid hotel ID", nil))
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, model.Error("Invalid image ID", nil))
		return
	}

	err = h.hotelImageUsecase.Delete(hotelID, imageID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[HANDLER][HotelImage][Delete] failed",
			"hotelID", hotelID,
			"imageID", imageID,
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
		"[HANDLER][HotelImage][Delete] success",
		"hotelID", hotelID,
		"imageID", imageID,
	)

	response.WriteJSON(w, http.StatusOK, model.Success("Image deleted", nil))
}
