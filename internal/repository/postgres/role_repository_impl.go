package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"strings"
)

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByName(name string) (*domain.Role, error) {
	query := `SELECT id, name FROM roles WHERE name = $1`

	var role domain.Role

	err := r.db.QueryRow(query, strings.ToLower(name)).
		Scan(&role.ID, &role.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.ErrorLogger.Println("[ROLE_REPO][GetByName] error:", err)
		return nil, err
	}

	return &role, nil
}
