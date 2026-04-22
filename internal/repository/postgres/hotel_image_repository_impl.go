package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"

	"github.com/google/uuid"
)

type hotelImageRepository struct {
	db *sql.DB
}

func NewHotelImageRepository(db *sql.DB) *hotelImageRepository {
	return &hotelImageRepository{db: db}
}

func (r *hotelImageRepository) Create(img *domain.HotelImage) error {
	query := `
		INSERT INTO hotel_images (
			id, hotel_id, image_url, is_primary, sort_order
		)
		VALUES ($1,$2,$3,$4,$5)
	`
	logger.InfoLogger.Println("[REPO][HotelImage][Create] inserting hotel image id:", img.ID)

	_, err := r.db.Exec(
		query,
		img.ID,
		img.HotelID,
		img.ImageURL,
		img.IsPrimary,
		img.SortOrder,
	)

	if err != nil {
		logger.ErrorLogger.Println(
			"[REPO][HotelImage][Create] failed insert hotel image id:",
			img.ID,
			"error:", err,
		)
		return err
	}

	logger.InfoLogger.Println("[REPO][HotelImage][Create] success insert hotel image id:", img.ID)

	return nil
}

func (r *hotelImageRepository) CountByHotelID(hotelID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM hotel_images 
		WHERE hotel_id = $1
		AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(query, hotelID).Scan(&count)
	if err != nil {
		logger.ErrorLogger.Println("[REPO][HotelImage][Count] error:", err)
		return 0, err
	}

	logger.InfoLogger.Println("[REPO][HotelImage][Create] success count:", count)

	return count, nil
}

func (r *hotelImageRepository) SetAllNonPrimary(hotelID uuid.UUID) error {
	query := `
			UPDATE hotel_images 
			SET is_primary = false 
			WHERE hotel_id = $1
			AND deleted_at IS NULL
		`
	_, err := r.db.Exec(query, hotelID)
	return err
}

func (r *hotelImageRepository) SetPrimary(imageID uuid.UUID) error {
	query := `
			UPDATE hotel_images 
			SET is_primary = true 
			WHERE id = $1
			AND deleted_at IS NULL
		`
	_, err := r.db.Exec(query, imageID)
	return err
}

func (r *hotelImageRepository) Exists(imageID, hotelID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM hotel_images 
			WHERE id = $1 
			AND hotel_id = $2
			AND deleted_at IS NULL
		)
	`
	var exists bool
	err := r.db.QueryRow(query, imageID, hotelID).Scan(&exists)
	return exists, err
}

func (r *hotelImageRepository) GetByHotelID(hotelID uuid.UUID) ([]string, error) {

	rows, err := r.db.Query(
		`SELECT image_url 
		 FROM hotel_images 
		 WHERE hotel_id = $1 
		 AND deleted_at IS NULL`,
		hotelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		images = append(images, url)
	}

	return images, nil
}

func (r *hotelImageRepository) Delete(imageID uuid.UUID) error {
	query := `
	UPDATE hotel_images
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, imageID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
