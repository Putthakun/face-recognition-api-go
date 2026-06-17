package usecase

import (
	"errors"

	"github.com/Putthakun/face-recognition-api-go/internal/domain/repository"
	"github.com/Putthakun/face-recognition-api-go/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid employee ID or password")

type AuthUsecase interface {
	Login(empID int64, password string) (token string, expiresAt int64, err error)
}

type authUsecase struct {
	empRepo    repository.EmployeeRepository
	jwtService jwt.Service
}

func NewAuthUsecase(empRepo repository.EmployeeRepository, jwtService jwt.Service) AuthUsecase {
	return &authUsecase{empRepo: empRepo, jwtService: jwtService}
}

func (u *authUsecase) Login(empID int64, password string) (string, int64, error) {
	cred, err := u.empRepo.FindCredentialByEmpID(empID)
	if err != nil || cred == nil {
		return "", 0, ErrInvalidCredentials
	}

	if !cred.IsActive {
		return "", 0, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)); err != nil {
		return "", 0, ErrInvalidCredentials
	}

	token, expiresAt, err := u.jwtService.Generate(empID, cred.Role)
	if err != nil {
		return "", 0, err
	}

	return token, expiresAt, nil
}
