package usecase

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"mime/multipart"
	"time"

	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
	"github.com/Putthakun/face-recognition-api-go/internal/infrastructure/cache"
	"github.com/Putthakun/face-recognition-api-go/internal/infrastructure/httpclient"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrFaceNotDetected       = errors.New("no face detected in the uploaded photo")
)

type CreateEmployeeInput struct {
	EmpID    int64
	Name     string
	Password string
	Role     string
	Photo    *multipart.FileHeader
}

type UpdateEmployeeInput struct {
	Name     string
	Password string
	Role     string
	IsActive *bool
	Photo    *multipart.FileHeader
}

type EmployeeResponse struct {
	EmpID     int64      `json:"empId"`
	Name      string     `json:"name"`
	Role      *string    `json:"role"`
	IsActive  *bool      `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt"`
}

type EmployeeUsecase interface {
	Create(input CreateEmployeeInput) (*EmployeeResponse, error)
	GetAll() ([]EmployeeResponse, error)
	Update(empID int64, input UpdateEmployeeInput) (*EmployeeResponse, error)
	Delete(empID int64) error
}

type employeeUsecase struct {
	empRepo    repository.EmployeeRepository
	faceClient httpclient.FaceClient
	cache      cache.FaceVectorCache
}

func NewEmployeeUsecase(
	empRepo repository.EmployeeRepository,
	faceClient httpclient.FaceClient,
	cache cache.FaceVectorCache,
) EmployeeUsecase {
	return &employeeUsecase{empRepo: empRepo, faceClient: faceClient, cache: cache}
}

func (u *employeeUsecase) Create(input CreateEmployeeInput) (*EmployeeResponse, error) {
	exists, err := u.empRepo.ExistsByID(input.EmpID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmployeeAlreadyExists
	}

	// Extract face vector BEFORE any DB writes
	var vector []float32
	if input.Photo != nil {
		vector, err = u.faceClient.GetEmbedding(input.Photo)
		if err != nil {
			return nil, err
		}
		if vector == nil {
			return nil, ErrFaceNotDetected
		}
	}

	emp := &entity.Employee{EmpID: input.EmpID, Name: input.Name}

	var cred *entity.Credential
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		role := input.Role
		if role == "" {
			role = "Employee"
		}
		cred = &entity.Credential{
			EmpID:        input.EmpID,
			PasswordHash: string(hash),
			Role:         role,
		}
	}

	if err := u.empRepo.CreateWithCredential(emp, cred); err != nil {
		return nil, err
	}

	if vector != nil {
		emb := &entity.FaceEmbedded{
			EmpID:            input.EmpID,
			FaceEmbeddedData: floatToBytes(vector),
		}
		if err := u.empRepo.CreateEmbedding(emb); err != nil {
			return nil, err
		}
		_ = u.cache.Set(input.EmpID, vector)
		u.notifyReload()
	}

	return toResponse(emp, cred), nil
}

func (u *employeeUsecase) GetAll() ([]EmployeeResponse, error) {
	employees, err := u.empRepo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]EmployeeResponse, len(employees))
	for i, e := range employees {
		result[i] = *toResponse(&e, e.Credential)
	}
	return result, nil
}

func (u *employeeUsecase) Update(empID int64, input UpdateEmployeeInput) (*EmployeeResponse, error) {
	emp, err := u.empRepo.FindByID(empID)
	if err != nil || emp == nil {
		return nil, ErrEmployeeNotFound
	}

	// Extract face vector BEFORE any DB writes
	var vector []float32
	if input.Photo != nil {
		vector, err = u.faceClient.GetEmbedding(input.Photo)
		if err != nil {
			return nil, err
		}
		if vector == nil {
			return nil, ErrFaceNotDetected
		}
	}

	if input.Name != "" {
		emp.Name = input.Name
		if err := u.empRepo.Save(emp); err != nil {
			return nil, err
		}
	}

	cred := emp.Credential
	if input.Password != "" || input.Role != "" || input.IsActive != nil {
		if cred == nil {
			hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			role := input.Role
			if role == "" {
				role = "Employee"
			}
			cred = &entity.Credential{
				EmpID:        empID,
				PasswordHash: string(hash),
				Role:         role,
			}
			if err := u.empRepo.CreateCredential(cred); err != nil {
				return nil, err
			}
		} else {
			if input.Password != "" {
				hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
				cred.PasswordHash = string(hash)
			}
			if input.Role != "" {
				cred.Role = input.Role
			}
			if input.IsActive != nil {
				cred.IsActive = *input.IsActive
			}
			if err := u.empRepo.SaveCredential(cred); err != nil {
				return nil, err
			}
		}
	}

	if vector != nil {
		existing, _ := u.empRepo.FindEmbedding(empID)
		if existing != nil {
			existing.FaceEmbeddedData = floatToBytes(vector)
			existing.CreatedAt = time.Now().UTC()
			_ = u.empRepo.SaveEmbedding(existing)
		} else {
			_ = u.empRepo.CreateEmbedding(&entity.FaceEmbedded{
				EmpID:            empID,
				FaceEmbeddedData: floatToBytes(vector),
			})
		}
		_ = u.cache.Set(empID, vector)
		u.notifyReload()
	}

	return toResponse(emp, cred), nil
}

func (u *employeeUsecase) Delete(empID int64) error {
	emp, err := u.empRepo.FindByID(empID)
	if err != nil || emp == nil {
		return ErrEmployeeNotFound
	}
	if err := u.empRepo.Delete(emp); err != nil {
		return err
	}
	_ = u.cache.Remove(empID)
	u.notifyReload()
	return nil
}

func (u *employeeUsecase) notifyReload() {
	if err := u.faceClient.TriggerReload(); err != nil {
		log.Printf("warn: could not notify face server to reload: %v", err)
	}
}

func toResponse(emp *entity.Employee, cred *entity.Credential) *EmployeeResponse {
	r := &EmployeeResponse{EmpID: emp.EmpID, Name: emp.Name}
	if cred != nil {
		r.Role = &cred.Role
		r.IsActive = &cred.IsActive
		r.CreatedAt = &cred.CreatedAt
	}
	return r
}

func floatToBytes(v []float32) []byte {
	b := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(f))
	}
	return b
}
