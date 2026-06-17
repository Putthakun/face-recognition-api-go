package repository

import "github.com/Putthakun/face-recognition-api-go/internal/domain/entity"

type EmployeeRepository interface {
	ExistsByID(empID int64) (bool, error)
	Create(emp *entity.Employee) error
	CreateCredential(cred *entity.Credential) error
	FindAll() ([]entity.Employee, error)
	FindByID(empID int64) (*entity.Employee, error)
	Save(emp *entity.Employee) error
	SaveCredential(cred *entity.Credential) error
	Delete(emp *entity.Employee) error

	// Face embedding
	FindEmbedding(empID int64) (*entity.FaceEmbedded, error)
	CreateEmbedding(emb *entity.FaceEmbedded) error
	SaveEmbedding(emb *entity.FaceEmbedded) error
	AllEmbeddings() ([]entity.FaceEmbedded, error)

	// Credential for auth
	FindCredentialByEmpID(empID int64) (*entity.Credential, error)
	AnyCredential() (bool, error)
	CreateWithCredential(emp *entity.Employee, cred *entity.Credential) error
}
