package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"

	"github.com/google/uuid"
)

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *roomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(h *domain.Room) error {
	query := `
		INSERT INTO rooms (
			id, hotel_id, name, description, capacity, base_price
		)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	logger.InfoLogger.Println("[REPO][Room][Create] inserting room id:", h.ID)

	_, err := r.db.Exec(
		query,
		h.ID,
		h.HotelID,
		h.Name,
		h.Description,
		h.Capacity,
		h.BasePrice,
	)

	if err != nil {
		logger.ErrorLogger.Println(
			"[REPO][Room][Create] failed insert room id:",
			h.ID,
			"error:", err,
		)
		return err
	}

	logger.InfoLogger.Println("[REPO][Room][Create] success insert room id:", h.ID)

	return nil
}

func (r *roomRepository) Delete(roomID uuid.UUID) error {
	query := `
		UPDATE rooms
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, roomID)
	if err != nil {
		logger.ErrorLogger.Println("[REPO][Room][Delete] error:", err)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	logger.InfoLogger.Println("[REPO][Room][Delete] success:", roomID)

	return nil
}

func (r *roomRepository) ListByHotelID(hotelID uuid.UUID) ([]domain.Room, error) {

	query := `
		SELECT id, hotel_id, name, description, capacity, base_price
		FROM rooms
		WHERE hotel_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room

	for rows.Next() {
		var room domain.Room
		err := rows.Scan(
			&room.ID,
			&room.HotelID,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.BasePrice,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
