package dto

import "github.com/gofiber/fiber/v2"

type LoginUserRequestDto struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

func NewLoginUserRequestDTO(ctx *fiber.Ctx) (*LoginUserRequestDto, error) {
	var data LoginUserRequestDto
	err := ctx.BodyParser(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
