package postgres

import (
	"database/sql"
	"hotel-loyalty/internal/domain"
)

type TokenRepositoryImpl struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepositoryImpl {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) Save(token *domain.UserToken) error {
	query := `
		INSERT INTO user_tokens (id, user_id, refresh_token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.RefreshToken,
		token.UserAgent,
		token.IPAddress,
		token.ExpiresAt,
	)

	return err
}

func (r *TokenRepositoryImpl) FindByRefreshToken(rt string) (*domain.UserToken, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at
		FROM user_tokens
		WHERE refresh_token = $1
	`

	var token domain.UserToken

	err := r.db.QueryRow(query, rt).Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
		&token.UserAgent,
		&token.IPAddress,
		&token.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepositoryImpl) DeleteByRefreshToken(rt string) error {
	query := `DELETE FROM user_tokens WHERE refresh_token = $1`
	_, err := r.db.Exec(query, rt)
	return err
}
