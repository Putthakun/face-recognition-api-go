package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Generate(empID int64, role string) (token string, expiresAt int64, err error)
	Validate(tokenString string) (*Claims, error)
}

type Claims struct {
	EmpID int64  `json:"empId"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
}

func NewService(secret, issuer, audience string, expiresHours int) Service {
	return &jwtService{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
		ttl:      time.Duration(expiresHours) * time.Hour,
	}
}

func (s *jwtService) Generate(empID int64, role string) (string, int64, error) {
	expiresAt := time.Now().Add(s.ttl)
	claims := &Claims{
		EmpID: empID,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, expiresAt.Unix(), nil
}

func (s *jwtService) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithAudience(s.audience))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
