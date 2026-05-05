package room

import (
	"database/sql"
	stdErrors "errors"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/repository"
	"time"

	"github.com/google/uuid"
)

type RoomUsecase struct {
	roomRepo  repository.RoomRepository
	hotelRepo repository.HotelRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository, hotelRepo repository.HotelRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepo:  roomRepo,
		hotelRepo: hotelRepo,
	}
}

func (u *RoomUsecase) Create(room *domain.Room) (*response.CreateRoomResponse, error) {

	if room.Name == "" {
		logger.InfoLogger.Println("[USECASE][Room][Create] name is required")
		return nil, errors.ErrBadRequest
	}

	if room.HotelID == uuid.Nil {
		logger.InfoLogger.Println("[USECASE][Room][Create] invalid hotel id")
		return nil, errors.ErrBadRequest
	}

	if room.Capacity <= 0 {
		logger.InfoLogger.Println("[USECASE][Room][Create] capacity must be greater than 0")
		return nil, errors.ErrBadRequest
	}

	if room.BasePrice <= 0 {
		logger.InfoLogger.Println("[USECASE][Room][Create] Base price must be greater than 0")
		return nil, errors.ErrBadRequest
	}

	// 🔥 generate ID & timestamp
	room.ID = uuid.New()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	// 🔥 insert ke DB
	err := u.roomRepo.Create(room)
	if err != nil {
		logger.ErrorLogger.Println("[USECASE][Room][Create] create failed:", err)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][Room][Create] success",
		"room=", room.Name,
		"hotel=", room.HotelID,
	)

	// 🔥 mapping response
	return mapper.ToRoomResponse(room), nil
}

func (u *RoomUsecase) List(hotelID uuid.UUID) ([]response.RoomListItem, error) {
	data, err := u.roomRepo.ListByHotelID(hotelID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][Room][List] failed",
			"hotelID", hotelID,
			"error", err,
		)
		return nil, errors.ErrInternal
	}

	var res []response.RoomListItem
	for _, r := range data {
		res = append(res, mapper.ToRoomListItem(r))
	}

	logger.InfoLogger.Println(
		"[USECASE][Room][List] success",
		"hotelID", hotelID,
		"count", len(res),
	)

	return res, nil
}

func (u *RoomUsecase) Delete(id uuid.UUID) error {
	err := u.roomRepo.Delete(id)
	if err != nil {
		logger.ErrorLogger.Println("[USECASE][Room][Delete] failed:", err)

		if stdErrors.Is(err, sql.ErrNoRows) {
			return errors.NewBadRequest("Room not found")
		}

		return errors.ErrInternal
	}
	logger.InfoLogger.Println("[USECASE][Room][Delete] success", "id", id)

	return nil
}
