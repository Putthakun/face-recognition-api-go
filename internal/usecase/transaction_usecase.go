package usecase

import (
	"time"

	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
)

type TransactionResponse struct {
	TransactionID  int64      `json:"transactionId"`
	EmpID          *int64     `json:"empId"`
	EmpName        *string    `json:"empName"`
	CameraID       *int64     `json:"cameraId"`
	CameraLocation *string    `json:"cameraLocation"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type TransactionUsecase interface {
	Create(empID, cameraID *int64) (*entity.Transaction, error)
	GetPaginated(page, pageSize int, sortDesc bool) ([]TransactionResponse, int64, error)
}

type transactionUsecase struct {
	repo repository.TransactionRepository
}

func NewTransactionUsecase(repo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{repo: repo}
}

func (u *transactionUsecase) Create(empID, cameraID *int64) (*entity.Transaction, error) {
	tx := &entity.Transaction{EmpID: empID, CameraID: cameraID}
	if err := u.repo.Create(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (u *transactionUsecase) GetPaginated(page, pageSize int, sortDesc bool) ([]TransactionResponse, int64, error) {
	txs, total, err := u.repo.FindPaginated(page, pageSize, sortDesc)
	if err != nil {
		return nil, 0, err
	}

	result := make([]TransactionResponse, len(txs))
	for i, tx := range txs {
		r := TransactionResponse{
			TransactionID: tx.TransactionID,
			EmpID:         tx.EmpID,
			CameraID:      tx.CameraID,
			CreatedAt:     tx.CreatedAt,
		}
		if tx.Employee != nil {
			r.EmpName = &tx.Employee.Name
		}
		if tx.Camera != nil {
			r.CameraLocation = &tx.Camera.Location
		}
		result[i] = r
	}
	return result, total, nil
}
