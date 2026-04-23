package postgres

import (
	"database/sql"
	"fmt"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model/response"

	"github.com/google/uuid"
)

type hotelRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) *hotelRepository {
	return &hotelRepository{db: db}
}

func (r *hotelRepository) Create(h *domain.Hotel) error {
	query := `
		INSERT INTO hotels (
			id, name, description, address,
			city, country, lat, lng,
			star_rating, status, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`

	logger.InfoLogger.Println("[REPO][Hotel][Create] inserting hotel id:", h.ID)

	_, err := r.db.Exec(
		query,
		h.ID,
		h.Name,
		h.Description,
		h.Address,
		h.City,
		h.Country,
		h.Lat,
		h.Lng,
		h.StarRating,
		h.Status,
		h.CreatedAt,
		h.UpdatedAt,
	)

	if err != nil {
		logger.ErrorLogger.Println(
			"[REPO][Hotel][Create] failed insert hotel id:",
			h.ID,
			"error:", err,
		)
		return err
	}

	logger.InfoLogger.Println("[REPO][Hotel][Create] success insert hotel id:", h.ID)

	return nil
}

func (r *hotelRepository) List(city string, star int, limit, offset int) ([]response.HotelListItem, int, error) {

	// 🔥 main query
	query := `
		SELECT 
			h.id, 
			h.name, 
			h.description,        -- 🔥 TAMBAHIN
			h.address,            -- 🔥 TAMBAHIN
			h.city, 
			h.country, 
			h.star_rating,
			COALESCE(hi.image_url, '')
		FROM hotels h
		LEFT JOIN hotel_images hi 
			ON hi.hotel_id = h.id 
			AND hi.is_primary = true
			AND hi.deleted_at IS NULL
		WHERE h.status = 'published'
		AND h.deleted_at IS NULL
	`

	// 🔥 count query
	countQuery := `
	SELECT COUNT(*) 
	FROM hotels h
	WHERE h.status = 'published'
	AND h.deleted_at IS NULL
	`

	args := []interface{}{}
	countArgs := []interface{}{}
	idx := 1

	// 🔥 filter city
	if city != "" {
		query += fmt.Sprintf(" AND h.city ILIKE $%d", idx)
		countQuery += fmt.Sprintf(" AND h.city ILIKE $%d", idx)

		args = append(args, "%"+city+"%")
		countArgs = append(countArgs, "%"+city+"%")
		idx++
	}

	// 🔥 filter star
	if star != 0 {
		query += fmt.Sprintf(" AND h.star_rating = $%d", idx)
		countQuery += fmt.Sprintf(" AND h.star_rating = $%d", idx)

		args = append(args, star)
		countArgs = append(countArgs, star)
		idx++
	}

	// 🔥 sorting + pagination
	query += fmt.Sprintf(" ORDER BY h.created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	// 🔥 execute main query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hotels []response.HotelListItem

	for rows.Next() {
		var h response.HotelListItem
		err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Adresses, &h.City, &h.Country, &h.StarRating, &h.ImageURL)
		if err != nil {
			return nil, 0, err
		}
		hotels = append(hotels, h)
	}

	// 🔥 execute count query
	var total int
	err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return hotels, total, nil
}

func (r *hotelRepository) GetByID(id uuid.UUID) (*domain.Hotel, error) {

	query := `
	SELECT id, name, description, city, country, star_rating
	FROM hotels
	WHERE id = $1
	AND deleted_at IS NULL
	`

	var h domain.Hotel

	err := r.db.QueryRow(query, id).Scan(
		&h.ID,
		&h.Name,
		&h.Description,
		&h.City,
		&h.Country,
		&h.StarRating,
	)

	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (r *hotelRepository) Publish(id uuid.UUID) error {
	query := `
	UPDATE hotels
	SET status = 'published', updated_at = NOW()
	WHERE id = $1
	AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, id)
	return err
}

func (r *hotelRepository) Delete(id uuid.UUID) error {
	query := `
	UPDATE hotels
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
