package jwt

import (
	"errors"
	"hotel-loyalty/internal/infrastructure/logger"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte
var tokenExp time.Duration

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func InitJWT(secret string, exp int) {
	secretKey = []byte(secret)
	tokenExp = time.Duration(exp) * time.Second
}

func GenerateToken(userID, email, role string) (string, int, error) {
	expStr := os.Getenv("JWT_EXP")

	expSec, err := strconv.Atoi(expStr)
	if err != nil {
		logger.ErrorLogger.Println("[JWT][GenerateToken] invalid JWT_EXP, fallback to 3600:", err)
		expSec = 3600
	}

	tokenExp := time.Duration(expSec) * time.Second

	logger.InfoLogger.Println("[JWT][GenerateToken] generating token for userID:", userID, "role:", role)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(secretKey)
	if err != nil {
		logger.ErrorLogger.Println("[JWT][GenerateToken] failed signing token:", err)
		return "", 0, err
	}

	logger.InfoLogger.Println("[JWT][GenerateToken] token generated successfully for userID:", userID)

	return signed, expSec, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	logger.InfoLogger.Println("[JWT][ValidateToken] validating token")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.ErrorLogger.Println("[JWT][ValidateToken] invalid signing method")
			return nil, errors.New("invalid signing method")
		}

		return secretKey, nil
	})

	if err != nil {
		logger.ErrorLogger.Println("[JWT][ValidateToken] failed parsing token:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		logger.ErrorLogger.Println("[JWT][ValidateToken] invalid claims type")
		return nil, errors.New("invalid token claims")
	}

	if !token.Valid {
		logger.ErrorLogger.Println("[JWT][ValidateToken] token not valid")
		return nil, errors.New("invalid token")
	}

	logger.InfoLogger.Println("[JWT][ValidateToken] token valid for userID:", claims.UserID)

	return claims, nil
}
