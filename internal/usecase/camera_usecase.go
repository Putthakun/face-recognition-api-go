package usecase

import (
	"errors"

	"github.com/Putthakun/face-recognition-api-go/internal/domain/entity"
	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
)

var ErrCameraNotFound = errors.New("camera not found")

type CameraUsecase interface {
	GetAll() ([]entity.Camera, error)
	Create(location string) (*entity.Camera, error)
	Update(cameraID int64, location string) (*entity.Camera, error)
	Delete(cameraID int64) error
}

type cameraUsecase struct {
	repo repository.CameraRepository
}

func NewCameraUsecase(repo repository.CameraRepository) CameraUsecase {
	return &cameraUsecase{repo: repo}
}

func (u *cameraUsecase) GetAll() ([]entity.Camera, error) {
	return u.repo.FindAll()
}

func (u *cameraUsecase) Create(location string) (*entity.Camera, error) {
	cam := &entity.Camera{Location: location}
	if err := u.repo.Create(cam); err != nil {
		return nil, err
	}
	return cam, nil
}

func (u *cameraUsecase) Update(cameraID int64, location string) (*entity.Camera, error) {
	cam, err := u.repo.FindByID(cameraID)
	if err != nil || cam == nil {
		return nil, ErrCameraNotFound
	}
	cam.Location = location
	if err := u.repo.Save(cam); err != nil {
		return nil, err
	}
	return cam, nil
}

func (u *cameraUsecase) Delete(cameraID int64) error {
	cam, err := u.repo.FindByID(cameraID)
	if err != nil || cam == nil {
		return ErrCameraNotFound
	}
	return u.repo.Delete(cam)
}
