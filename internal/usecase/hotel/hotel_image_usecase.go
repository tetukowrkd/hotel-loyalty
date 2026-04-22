package hotel

import (
	stdErrors "errors"

	"database/sql"
	"fmt"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/pkg/validator"
	"hotel-loyalty/internal/repository"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type HotelImageUsecase struct {
	hotelImageRepository repository.HotelImageRepository
	db                   *sql.DB
}

func NewHotelImageUsecase(repo repository.HotelImageRepository, db *sql.DB) *HotelImageUsecase {
	return &HotelImageUsecase{
		hotelImageRepository: repo,
		db:                   db,
	}
}

func (u *HotelImageUsecase) Create(
	hotelID uuid.UUID,
	file multipart.File,
	filename string,
	size int64,
) (*response.CreateHotelImageResponse, error) {

	// 🔥 validate (ALL IN ONE)
	err := validator.ValidateImage(file, filename, size)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] validation failed:", err)
		return nil, err
	}

	// 🔥 generate filename
	ext := strings.ToLower(filepath.Ext(filename))
	newFileName := uuid.New().String() + ext

	dir := fmt.Sprintf("storage/hotels/%s", hotelID.String())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] mkdir failed:", err)
		return nil, errors.ErrInternal
	}

	fullPath := fmt.Sprintf("%s/%s", dir, newFileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] create file failed:", err)
		return nil, errors.ErrInternal
	}
	defer dst.Close()

	// 🔥 copy file
	if _, err := io.Copy(dst, file); err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] save file failed:", err)
		return nil, errors.ErrInternal
	}

	// 🔥 count existing images
	count, err := u.hotelImageRepository.CountByHotelID(hotelID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] count failed:", err)
		return nil, errors.ErrInternal
	}

	// 🔥 create domain
	img := &domain.HotelImage{
		ID:        uuid.New(),
		HotelID:   hotelID,
		ImageURL:  fmt.Sprintf("/storage/hotels/%s/%s", hotelID, newFileName),
		IsPrimary: count == 0,
		SortOrder: count + 1,
	}

	if err := u.hotelImageRepository.Create(img); err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] insert failed:", err)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println("[HOTEL_IMAGE_USECASE] success upload:", img.ImageURL)

	return mapper.ToHotelImageResponse(img), nil
}

func (u *HotelImageUsecase) SetPrimary(hotelID, imageID uuid.UUID) error {

	// 🔥 validasi image milik hotel
	exists, err := u.hotelImageRepository.Exists(imageID, hotelID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] exists check failed:", err)
		return errors.ErrInternal
	}
	if !exists {
		return errors.NewBadRequest("Image not found for this hotel")
	}

	// 🔥 begin transaction
	tx, err := u.db.Begin()
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] begin tx failed:", err)
		return errors.ErrInternal
	}
	defer tx.Rollback()

	// 🔥 reset primary
	_, err = tx.Exec(`UPDATE hotel_images SET is_primary = false WHERE hotel_id = $1`, hotelID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] reset primary failed:", err)
		return errors.ErrInternal
	}

	// 🔥 set new primary
	_, err = tx.Exec(`UPDATE hotel_images SET is_primary = true WHERE id = $1`, imageID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] set primary failed:", err)
		return errors.ErrInternal
	}

	// 🔥 commit
	if err := tx.Commit(); err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] commit failed:", err)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[HOTEL_IMAGE_USECASE] primary updated",
		"hotelID=", hotelID,
		"imageID=", imageID,
	)
	return nil
}

func (u *HotelImageUsecase) Delete(hotelID, imageID uuid.UUID) error {

	// 🔥 cek exists (repo harus pakai deleted_at IS NULL)
	exists, err := u.hotelImageRepository.Exists(imageID, hotelID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] exists check failed:", err)
		return errors.ErrInternal
	}
	if !exists {
		return errors.NewBadRequest("Image not found for this hotel")
	}

	// 🔥 begin tx
	tx, err := u.db.Begin()
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] begin tx failed:", err)
		return errors.ErrInternal
	}
	defer tx.Rollback()

	// 🔥 cek apakah primary (pakai hotel_id + deleted_at)
	var isPrimary bool
	err = tx.QueryRow(`
		SELECT is_primary 
		FROM hotel_images 
		WHERE id = $1 
		AND hotel_id = $2
		AND deleted_at IS NULL
	`, imageID, hotelID).Scan(&isPrimary)

	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return errors.NewBadRequest("Image not found")
		}
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] check primary failed:", err)
		return errors.ErrInternal
	}

	// 🔥 SOFT DELETE DULU
	res, err := tx.Exec(`
		UPDATE hotel_images
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, imageID)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] delete failed:", err)
		return errors.ErrInternal
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.NewBadRequest("Image already deleted or not found")
	}

	// 🔥 HANDLE PRIMARY (INI TEMPAT BLOK LO MASUK)
	if isPrimary {

		// 🔥 MATIKAN semua primary dulu
		_, err = tx.Exec(`
			UPDATE hotel_images
			SET is_primary = false
			WHERE hotel_id = $1
		`, hotelID)
		if err != nil {
			logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] reset primary failed:", err)
			return errors.ErrInternal
		}

		var nextID uuid.UUID

		// 🔥 cari pengganti
		err = tx.QueryRow(`
			SELECT id 
			FROM hotel_images
			WHERE hotel_id = $1
			AND deleted_at IS NULL
			AND id != $2
			ORDER BY sort_order ASC
			LIMIT 1
		`, hotelID, imageID).Scan(&nextID)

		if err != nil && !stdErrors.Is(err, sql.ErrNoRows) {
			logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] get next image failed:", err)
			return errors.ErrInternal
		}

		// 🔥 kalau ada → set jadi primary
		if err == nil {
			_, err = tx.Exec(`
				UPDATE hotel_images 
				SET is_primary = true 
				WHERE id = $1
			`, nextID)
			if err != nil {
				logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] set new primary failed:", err)
				return errors.ErrInternal
			}
		}
	}

	// 🔥 commit
	if err := tx.Commit(); err != nil {
		logger.ErrorLogger.Println("[HOTEL_IMAGE_USECASE] commit failed:", err)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[HOTEL_IMAGE_USECASE] delete success",
		"hotelID", hotelID,
		"imageID", imageID,
	)

	return nil
}
