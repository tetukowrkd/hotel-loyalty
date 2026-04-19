package identity

import (
	"time"

	"github.com/google/uuid"

	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/pkg/hash"
	"hotel-loyalty/internal/repository"
)

type TokenUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository

	refreshTokenExp time.Duration
}

func NewTokenUsecase(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, refreshExp int) *TokenUsecase {
	return &TokenUsecase{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		refreshTokenExp: time.Duration(refreshExp) * time.Second,
	}
}

func (u *TokenUsecase) RefreshToken(refreshToken string) (*response.RefreshResponse, error) {
	hashedToken := hash.HashToken(refreshToken)

	tokenData, err := u.tokenRepo.FindByRefreshToken(hashedToken)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if tokenData == nil {
		return nil, errors.ErrBadRequest
	}

	if time.Now().After(tokenData.ExpiresAt) {
		return nil, errors.ErrBadRequest
	}

	user, err := u.userRepo.GetByID(tokenData.UserID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrBadRequest
	}

	// 🔐 access token baru (WITH ROLE)
	newAccessToken, accessTTL, err := jwt.GenerateToken(
		user.ID.String(),
		user.Email,
		user.RoleName,
	)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// 🔁 refresh token baru
	newRefreshToken := uuid.New().String()
	hashedNew := hash.HashToken(newRefreshToken)

	// delete old
	_ = u.tokenRepo.DeleteByRefreshToken(hashedToken)

	// save new
	token := &domain.UserToken{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: hashedNew,
		ExpiresAt:    time.Now().Add(u.refreshTokenExp),
	}

	if err := u.tokenRepo.Save(token); err != nil {
		return nil, errors.ErrInternal
	}

	return &response.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken, // 🔥 FIXED
		ExpiresIn:    accessTTL,
	}, nil
}
