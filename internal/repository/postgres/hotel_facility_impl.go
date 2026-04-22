package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"

	"github.com/google/uuid"
)

type hotelFacilityRepository struct {
	db *sql.DB
}

func NewHotelFacilityRepository(db *sql.DB) *hotelFacilityRepository {
	return &hotelFacilityRepository{db: db}
}

func (r *hotelFacilityRepository) GetAll() ([]domain.HotelFacility, error) {
	query := `
			SELECT id, name, icon 
			FROM hotel_facilities
			WHERE deleted_at IS NULL
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facilities []domain.HotelFacility

	for rows.Next() {
		var f domain.HotelFacility
		if err := rows.Scan(&f.ID, &f.Name, &f.Icon); err != nil {
			return nil, err
		}
		facilities = append(facilities, f)
	}

	return facilities, nil
}

func (r *hotelFacilityRepository) Create(name, icon string) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.db.Exec(
		`INSERT INTO hotel_facilities (id, name, icon) VALUES ($1,$2,$3)`,
		id, name, icon,
	)

	return id, err
}

func (r *hotelFacilityRepository) Delete(id uuid.UUID) error {
	query := `
	UPDATE hotel_facilities
	SET deleted_at = NOW()
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

type hotelFacilityMapRepository struct {
	db *sql.DB
}

func NewHotelFacilityMapRepository(db *sql.DB) *hotelFacilityMapRepository {
	return &hotelFacilityMapRepository{db: db}
}

func (r *hotelFacilityMapRepository) Assign(hotelID uuid.UUID, facilityIDs []uuid.UUID) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 🔥 clear existing
	_, err = tx.Exec(`DELETE FROM hotel_facility_maps WHERE hotel_id = $1`, hotelID)
	if err != nil {
		return err
	}

	// 🔥 insert new
	for _, fid := range facilityIDs {
		_, err := tx.Exec(
			`INSERT INTO hotel_facility_maps (hotel_id, facility_id) VALUES ($1,$2)`,
			hotelID, fid,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *hotelFacilityMapRepository) GetByHotelID(hotelID uuid.UUID) ([]domain.HotelFacility, error) {

	query := `
		SELECT hf.id, hf.name, hf.icon
		FROM hotel_facilities hf
		JOIN hotel_facility_maps hfm 
		  ON hf.id = hfm.facility_id
		WHERE hfm.hotel_id = $1
		AND hf.deleted_at IS NULL
		AND hfm.deleted_at IS NULL
	`

	rows, err := r.db.Query(query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facilities []domain.HotelFacility

	for rows.Next() {
		var f domain.HotelFacility
		err := rows.Scan(&f.ID, &f.Name, &f.Icon)
		if err != nil {
			return nil, err
		}
		facilities = append(facilities, f)
	}

	return facilities, nil
}
