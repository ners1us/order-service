package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/ners1us/order-service/internal/model"
	"time"
)

type JWTService interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(tokenString string) (*model.Claims, error)
}

type jwtServiceImpl struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtServiceImpl{secretKey}
}

func (js *jwtServiceImpl) GenerateToken(userID, role string) (string, error) {
	claims := model.Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(js.secretKey))
}

func (js *jwtServiceImpl) ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
