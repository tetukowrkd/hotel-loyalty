package config

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	BaseUrl         string
	JWTSecret       string
	JWTExp          int
	RefreshTokenExp int
	Origins         string
}
