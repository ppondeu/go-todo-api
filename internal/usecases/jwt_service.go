package usecases

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
)

type JWTService interface {
	SignToken(claims *dtos.UserClaims, secret []byte) (*string, error)
	Validate(tokenString string, secret []byte) (*dtos.UserClaims, error)
	GetAccess() []byte
	GetRefresh() []byte
}

type jwtServiceImpl struct {
	jwtSecret     []byte
	refreshSecret []byte
}

func NewJwtService(jwtSecret, refreshSecret []byte) JWTService {
	return &jwtServiceImpl{jwtSecret: jwtSecret, refreshSecret: refreshSecret}
}

func (s *jwtServiceImpl) SignToken(claims *dtos.UserClaims, secret []byte) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *jwtServiceImpl) Validate(tokenString string, secret []byte) (*dtos.UserClaims, error) {
	claims := &dtos.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewBadRequestError("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		logs.Error(err)
		return nil, errs.NewBadRequestError("invalid refresh token")
	}

	return claims, nil
}

func (s *jwtServiceImpl) GetAccess() []byte {
	return s.jwtSecret
}

func (s *jwtServiceImpl) GetRefresh() []byte {
	return s.refreshSecret
}
