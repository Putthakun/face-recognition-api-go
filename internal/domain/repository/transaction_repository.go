package repository

import "github.com/Putthakun/face-recognition-api-go/internal/domain/entity"

type TransactionRepository interface {
	Create(tx *entity.Transaction) error
	FindPaginated(page, pageSize int, sortDesc bool) ([]entity.Transaction, int64, error)
}
