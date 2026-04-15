package usecase

import (
	"e-commerce/internal/domain"
	"e-commerce/internal/infrastructure/logger"
	"e-commerce/internal/pkg/errors"
	"e-commerce/internal/pkg/hash"
	"e-commerce/internal/pkg/jwt"
	"e-commerce/internal/repository"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
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

func (u *UserUsecase) Login(email, password string) (string, error) {
	// 1. ambil user
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Login] GetByEmail:", err)
		return "", errors.ErrInternal
	}

	if user == nil {
		return "", errors.ErrBadRequest
	}

	// 2. compare password
	if !hash.CheckPassword(password, user.Password) {
		return "", errors.ErrBadRequest
	}

	// 3. generate token
	token, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", errors.ErrInternal
	}

	return token, nil
}
