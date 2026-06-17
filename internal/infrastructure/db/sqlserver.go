package db

import (
	"fmt"

	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/pkg/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLServer(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable&TrustServerCertificate=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to SQL Server: %w", err)
	}

	// Auto-migrate — creates tables if they don't exist
	if err := db.AutoMigrate(
		&entity.Employee{},
		&entity.Credential{},
		&entity.FaceEmbedded{},
		&entity.Camera{},
		&entity.Transaction{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	return db, nil
}
