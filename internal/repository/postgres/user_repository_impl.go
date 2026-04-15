package postgres

import (
	"database/sql"
	"e-commerce/internal/domain"
	"e-commerce/internal/infrastructure/logger"
	"time"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		logger.ErrorLogger.Println("failed insert user:", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(query, email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
			return nil, nil
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return &user, nil
}
