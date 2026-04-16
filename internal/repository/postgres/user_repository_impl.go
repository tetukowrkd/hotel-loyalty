package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"time"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

const baseUserQuery = `
	SELECT 
		u.id,
		u.name,
		u.email,
		u.password,
		u.role_id,
		r.name AS role_name,
		u.phone,
		u.member_code,
		u.is_active,
		u.created_at,
		u.updated_at
	FROM users u
	LEFT JOIN roles r ON u.role_id = r.id
`

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			id, name, email, password, role_id, phone, member_code,
			is_active, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
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
		user.RoleID,
		user.Phone,
		user.MemberCode,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		logger.ErrorLogger.Println("[USER_REPO][Create] error:", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := baseUserQuery + `
		WHERE u.email = $1
	`

	row := r.db.QueryRow(query, email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.RoleName,
		&user.Phone,
		&user.MemberCode,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.ErrorLogger.Println("[USER_REPO][GetByEmail] error:", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	query := baseUserQuery + `
		WHERE u.id = $1
	`

	row := r.db.QueryRow(query, id)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.RoleName,
		&user.Phone,
		&user.MemberCode,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.ErrorLogger.Println("[USER_REPO][GetByID] error:", err)
		return nil, err
	}

	return &user, nil
}
