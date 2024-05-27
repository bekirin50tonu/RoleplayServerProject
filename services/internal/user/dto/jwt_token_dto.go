package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTTokenDTO struct {
	Username string             `json:"username"`
	ID       primitive.ObjectID `json:"id"`
	jwt.RegisteredClaims
}

func NewJWTTokenDTO(username string, id primitive.ObjectID, exp time.Time) JWTTokenDTO {

	return JWTTokenDTO{
		Username: username,
		ID:       id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

}
