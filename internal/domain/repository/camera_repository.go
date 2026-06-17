package repository

import "github.com/Putthakun/face-recognition-api-go/internal/domain/entity"

type CameraRepository interface {
	FindAll() ([]entity.Camera, error)
	FindByID(cameraID int64) (*entity.Camera, error)
	Create(camera *entity.Camera) error
	Save(camera *entity.Camera) error
	Delete(camera *entity.Camera) error
}
