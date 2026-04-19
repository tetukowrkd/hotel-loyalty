package identity

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/pkg/generator"
	"hotel-loyalty/internal/pkg/hash"
	"hotel-loyalty/internal/pkg/qrcode"
	"hotel-loyalty/internal/repository"
)

type UserUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository

	refreshTokenExp time.Duration
}

func NewUserUsecase(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, refreshExp int) *UserUsecase {
	return &UserUsecase{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		refreshTokenExp: time.Duration(refreshExp) * time.Second,
	}
}

func (u *UserUsecase) Register(user *domain.User) (*response.RegisterUserResponse, error) {
	existingUser, err := u.userRepo.GetByEmail(user.Email)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Register][GetByEmail] error:", err)
		return nil, errors.ErrInternal
	}

	if existingUser != nil {
		logger.InfoLogger.Println("[USER_USECASE][Register] email exists:", user.Email)
		return nil, errors.ErrEmailExists
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][Register][HashPassword] error:", err)
		return nil, errors.ErrInternal
	}

	user.Password = hashedPassword
	user.IsActive = true

	// 🔥 default role
	memberRoleID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	user.RoleID = &memberRoleID
	user.RoleName = "member"

	var createErr error

	for i := 0; i < 3; i++ {
		user.MemberCode = generator.GenerateMemberCode(user.RoleName)

		createErr = u.userRepo.Create(user)
		if createErr == nil {
			logger.InfoLogger.Println(
				"[USER_USECASE][Register] success",
				"email=", user.Email,
				"member_code=", user.MemberCode,
				"attempt=", i+1,
			)

			return mapper.ToRegisterResponse(user), nil
		}

		logger.ErrorLogger.Println(
			"[USER_USECASE][Register] create failed",
			"attempt=", i+1,
			"error=", createErr,
		)
	}

	// 🔥 kalau semua retry gagal
	logger.ErrorLogger.Println(
		"[USER_USECASE][Register] failed after retries",
		"email=", user.Email,
		"error=", createErr,
	)

	return nil, errors.ErrInternal
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
		ExpiresAt:    time.Now().Add(u.refreshTokenExp),
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

func (u *UserUsecase) GetMemberQR(userID string) (string, error) {
	logger.InfoLogger.Println("[USER_USECASE][GetMemberQR] start userID:", userID)

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][GetMemberQR][GetByID] error:", err)
		return "", errors.ErrInternal
	}

	if user == nil {
		logger.InfoLogger.Println("[USER_USECASE][GetMemberQR] user not found:", userID)
		return "", errors.ErrNotFound
	}

	if user.MemberCode == "" {
		logger.ErrorLogger.Println("[USER_USECASE][GetMemberQR] empty member_code user:", userID)
		return "", errors.ErrInternal
	}

	png, err := qrcode.GenerateQRCode(user.MemberCode)
	if err != nil {
		logger.ErrorLogger.Println("[USER_USECASE][GetMemberQR][GenerateQR] error:", err)
		return "", errors.ErrInternal
	}

	qrBase64 := base64.StdEncoding.EncodeToString(png)

	logger.InfoLogger.Println("[USER_USECASE][GetMemberQR] success userID:", userID)

	return qrBase64, nil
}
