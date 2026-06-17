package config

import (
	"os"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisAddr     string
	RedisPassword string

	// JWT
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	JWTExpiresHr int

	// Face server
	FaceAPIBaseURL string

	// Admin seed
	AdminEmpID   int64
	AdminPassword string

	// CORS
	AllowedOrigins []string

	// Server
	Port string
}

func Load() *Config {
	cfg := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "1433"),
		DBUser:         getEnv("DB_USER", "sa"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "facerecog"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-at-least-32-chars-long!!"),
		JWTIssuer:      getEnv("JWT_ISSUER", "face-recognition-api-go"),
		JWTAudience:    getEnv("JWT_AUDIENCE", "face-recognition-client"),
		JWTExpiresHr:   24,
		FaceAPIBaseURL: getEnv("FACE_API_BASE_URL", "http://localhost:8001"),
		AdminEmpID:     1111,
		AdminPassword:  getEnv("ADMIN_PASSWORD", "Admin@1234"),
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: []string{getEnv("ALLOWED_ORIGIN", "http://localhost:5173")},
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
