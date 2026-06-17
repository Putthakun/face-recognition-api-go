package repository

import (
	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type transactionRepo struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(tx *entity.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepo) FindPaginated(page, pageSize int, sortDesc bool) ([]entity.Transaction, int64, error) {
	var total int64
	if err := r.db.Model(&entity.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Preload("Employee").Preload("Camera")
	if sortDesc {
		query = query.Order("CreatedAt DESC")
	} else {
		query = query.Order("CreatedAt ASC")
	}

	var txs []entity.Transaction
	err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&txs).Error

	return txs, total, err
}
