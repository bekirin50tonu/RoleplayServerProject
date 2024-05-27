package dto

import "github.com/gofiber/fiber/v2"

type RegisterUserRequestDTO struct {
	Name     string `validate:"required" json:"name"`
	LastName string `validate:"required" json:"lastname"`
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

func NewRegisterUserRequestDTO(ctx *fiber.Ctx) (*RegisterUserRequestDTO, error) {
	var data RegisterUserRequestDTO
	err := ctx.BodyParser(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
