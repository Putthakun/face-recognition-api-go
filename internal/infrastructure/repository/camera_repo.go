package repository

import (
	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

type cameraRepo struct{ db *gorm.DB }

func NewCameraRepository(db *gorm.DB) repository.CameraRepository {
	return &cameraRepo{db: db}
}

func (r *cameraRepo) FindAll() ([]entity.Camera, error) {
	var cameras []entity.Camera
	err := r.db.Find(&cameras).Error
	return cameras, err
}

func (r *cameraRepo) FindByID(cameraID int64) (*entity.Camera, error) {
	var cam entity.Camera
	err := r.db.First(&cam, "CameraId = ?", cameraID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &cam, err
}

func (r *cameraRepo) Create(camera *entity.Camera) error {
	return r.db.Create(camera).Error
}

func (r *cameraRepo) Save(camera *entity.Camera) error {
	return r.db.Save(camera).Error
}

func (r *cameraRepo) Delete(camera *entity.Camera) error {
	return r.db.Delete(camera).Error
}
