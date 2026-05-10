package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"time"

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

func (r *roomRepository) UpsertInventory(roomID uuid.UUID, date time.Time, stock int) error {
	query := `
		INSERT INTO room_inventory (id, room_id, date, available_stock)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (room_id, date)
		DO UPDATE SET available_stock = EXCLUDED.available_stock
	`

	_, err := r.db.Exec(
		query,
		uuid.New(),
		roomID,
		date,
		stock,
	)

	if err != nil {
		logger.ErrorLogger.Println("[REPO][Inventory][Upsert] failed:", err)
		return err
	}

	return nil
}

func (r *roomRepository) Exists(roomID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM rooms
			WHERE id = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.QueryRow(query, roomID).Scan(&exists)
	if err != nil {
		logger.ErrorLogger.Println("[REPO][Room][Exists] error:", err)
		return false, err
	}

	return exists, nil
}
