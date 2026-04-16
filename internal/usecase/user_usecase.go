package usecase

import (
	"time"

	"github.com/google/uuid"

	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
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

func (u *UserUsecase) Register(user *domain.User) (*response.RegisterUserResponse, error) {
	existingUser, err := u.userRepo.GetByEmail(user.Email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][GetByEmail] Error:", err)
		return nil, errors.ErrInternal
	}

	if existingUser != nil {
		return nil, errors.ErrEmailExists
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	user.Password = hashedPassword
	user.IsActive = true

	// 🔥 SET DEFAULT ROLE (USER)
	memberRoleID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	user.RoleID = &memberRoleID

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.ErrInternal
	}

	return mapper.ToRegisterResponse(user), nil
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

	accessToken, accessTTL, err := jwt.GenerateToken(user.ID.String(), user.Email, user.RoleName)
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
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
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
