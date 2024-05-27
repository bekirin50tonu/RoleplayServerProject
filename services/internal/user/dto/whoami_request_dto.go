package dto

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type WhoAmIRequestDto struct {
	AccessToken string `json:"X-Token"`
}

func NewWhoAmIRequestDTO(ctx *fiber.Ctx) (*WhoAmIRequestDto, error) {
	token := ctx.Get("X-Token")
	fmt.Print("token", token)
	resp := WhoAmIRequestDto{AccessToken: token}
	return &resp, nil
}
