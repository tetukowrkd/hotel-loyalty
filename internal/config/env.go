package config

import (
	"fmt"
	"hotel-loyalty/internal/infrastructure/logger"
	"os"
	"strconv"
)

func LoadConfig() *Config {
	jwtExp, err := strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		logger.InfoLogger.Println("[Config][ENV] invalid JWT_EXP")
		panic("Invalid JWT_EXP")
	}

	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		logger.InfoLogger.Println("[Config][ENV] invalid REFRESH_TOKEN_EXP")
		panic("Invalid REFRESH_TOKEN_EXP")
	}

	return &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExp:          jwtExp,
		RefreshTokenExp: refreshExp,
	}
}

func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWTExp == 0 {
		return fmt.Errorf("JWT_EXP invalid")
	}
	if c.RefreshTokenExp == 0 {
		return fmt.Errorf("REFRESH_TOKEN_EXP invalid")
	}
	return nil
}
