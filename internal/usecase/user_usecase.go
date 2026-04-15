package usecase

import (
	"time"

	"github.com/google/uuid"

	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/pkg/hash"
	"hotel-loyalty/internal/repository"
)

type UserUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewUserUsecase(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (u *UserUsecase) Register(user *domain.User) error {
	// 1. cek email sudah ada atau belum
	existingUser, err := u.userRepo.GetByEmail(user.Email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][GetByEmail] Error:", err)
		return errors.ErrInternal
	}

	if existingUser != nil {
		logger.InfoLogger.Println("[USER_USECASE][Register] Email Already Exists:", user.Email)
		return errors.ErrEmailExists
	}

	// 2. hash password
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][HashPassword] Error:", err)
		return errors.ErrInternal
	}

	user.Password = hashedPassword
	user.IsActive = true

	// 3. insert ke DB
	err = u.userRepo.Create(user)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][CreateUser] Error:", err)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println("[USER_USECASE][Register] User Created:", user.Email)

	return nil
}

func (u *UserUsecase) Login(email, password string) (*response.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Login][GetByEmail] error:", err)
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrBadRequest
	}

	if !hash.CheckPassword(password, user.Password) {
		logger.InfoLogger.Println("[USER_USECASE][Login] invalid credentials:", email)
		return nil, errors.ErrBadRequest
	}

	accessToken, accessTTL, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.ErrInternal
	}

	refreshToken := uuid.New().String()
	hashedToken := hash.HashToken(refreshToken)

	token := &domain.UserToken{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: hashedToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	err = u.tokenRepo.Save(token)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Login][SaveToken] error:", err)
		return nil, errors.ErrInternal
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTTL,
	}, nil
}

func (u *UserUsecase) Logout(refreshToken string) error {
	hashedToken := hash.HashToken(refreshToken)

	err := u.tokenRepo.DeleteByRefreshToken(hashedToken)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Logout][DeleteToken] error:", err)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println("[USER_USECASE][Logout] success")

	return nil
}

func (u *UserUsecase) GetProfile(userID string) (*domain.User, error) {
	return u.userRepo.GetByID(userID)
}

func (u *UserUsecase) RefreshToken(refreshToken string) (*response.RefreshResponse, error) {
	hashedToken := hash.HashToken(refreshToken)

	tokenData, err := u.tokenRepo.FindByRefreshToken(hashedToken)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][RefreshToken] Find:", err)
		return nil, errors.ErrInternal
	}

	if tokenData == nil {
		return nil, errors.ErrInternal
	}

	if time.Now().After(tokenData.ExpiresAt) {
		logger.InfoLogger.Println("[USER_USECASE][RefreshToken] expired token:", tokenData.UserID)
		return nil, errors.ErrInternal
	}

	user, err := u.userRepo.GetByID(tokenData.UserID.String())
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][RefreshToken] GetUser:", err)
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrInternal
	}

	// 🔐 new access token
	newAccessToken, accessTTL, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][RefreshToken] GenerateToken:", err)
		return nil, errors.ErrInternal
	}

	// 🔁 new refresh token
	newRefreshToken := uuid.New().String()
	hashedNew := hash.HashToken(newRefreshToken)

	// delete old
	if err := u.tokenRepo.DeleteByRefreshToken(hashedToken); err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][RefreshToken] Delete:", err)
	}

	// save new
	token := &domain.UserToken{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: hashedNew,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := u.tokenRepo.Save(token); err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][RefreshToken] Save:", err)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println("[USER_USECASE][RefreshToken] success:", user.Email)

	return &response.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTTL,
	}, nil
}
