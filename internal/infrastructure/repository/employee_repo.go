package repository

import (
	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type employeeRepo struct{ db *gorm.DB }

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepo{db: db}
}

func (r *employeeRepo) ExistsByID(empID int64) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Employee{}).Where("EmpId = ?", empID).Count(&count).Error
	return count > 0, err
}

func (r *employeeRepo) Create(emp *entity.Employee) error {
	return r.db.Create(emp).Error
}

func (r *employeeRepo) CreateCredential(cred *entity.Credential) error {
	return r.db.Create(cred).Error
}

func (r *employeeRepo) CreateWithCredential(emp *entity.Employee, cred *entity.Credential) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(emp).Error; err != nil {
			return err
		}
		if cred != nil {
			cred.EmpID = emp.EmpID
			if err := tx.Create(cred).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *employeeRepo) FindAll() ([]entity.Employee, error) {
	var employees []entity.Employee
	err := r.db.Preload("Credential").Find(&employees).Error
	return employees, err
}

func (r *employeeRepo) FindByID(empID int64) (*entity.Employee, error) {
	var emp entity.Employee
	err := r.db.Preload("Credential").First(&emp, "EmpId = ?", empID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &emp, err
}

func (r *employeeRepo) Save(emp *entity.Employee) error {
	return r.db.Save(emp).Error
}

func (r *employeeRepo) SaveCredential(cred *entity.Credential) error {
	return r.db.Save(cred).Error
}

func (r *employeeRepo) Delete(emp *entity.Employee) error {
	return r.db.Select("Credential", "FaceEmbeddeds", "Transactions").Delete(emp).Error
}

func (r *employeeRepo) FindEmbedding(empID int64) (*entity.FaceEmbedded, error) {
	var emb entity.FaceEmbedded
	err := r.db.Where("EmpId = ?", empID).First(&emb).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &emb, err
}

func (r *employeeRepo) CreateEmbedding(emb *entity.FaceEmbedded) error {
	return r.db.Create(emb).Error
}

func (r *employeeRepo) SaveEmbedding(emb *entity.FaceEmbedded) error {
	return r.db.Save(emb).Error
}

func (r *employeeRepo) AllEmbeddings() ([]entity.FaceEmbedded, error) {
	var embeddings []entity.FaceEmbedded
	err := r.db.Where("FaceEmbeddedData IS NOT NULL").Find(&embeddings).Error
	return embeddings, err
}

func (r *employeeRepo) FindCredentialByEmpID(empID int64) (*entity.Credential, error) {
	var cred entity.Credential
	err := r.db.Where("EmpId = ? AND IsActive = 1", empID).First(&cred).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &cred, err
}

func (r *employeeRepo) AnyCredential() (bool, error) {
	var count int64
	err := r.db.Model(&entity.Credential{}).Count(&count).Error
	return count > 0, err
}
