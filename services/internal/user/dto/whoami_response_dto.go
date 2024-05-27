package dto

import (
	"services/internal/user/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WhoAmIResponseDto struct {
	Username  string             `json:"username"`
	ID        primitive.ObjectID `json:"user_id"`
	Name      string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func NewWhoAmIResponseDTO(account models.Account) WhoAmIResponseDto {
	return WhoAmIResponseDto{
		Username:  account.Username,
		ID:        account.User.ID,
		Name: account.User.Name,
		LastName: account.User.Lastname,
		Email: account.User.Email,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
