package controllers

import (
	"services/internal/storage/dto"
	"services/internal/storage/services"
	"services/pkg/common/response"

	"github.com/gofiber/fiber/v2"
)

type StorageController struct {
	service_storage *services.StorageService
}

func NewStorageController(service *services.StorageService) *StorageController {
	return &StorageController{
		service_storage: service,
	}
}

func (c *StorageController) SaveFileUpdatedDataHandler(ctx *fiber.Ctx) error {

	diskName := ctx.Params("disk")
	form, err := ctx.MultipartForm()

	if err != nil {
		return err
	}

	file := form.File["images"][0]
	img, err := c.service_storage.SaveImageToLocalStorageRest(diskName, file)
	if err != nil {
		return err
	}
	resp := dto.NewUploadResponseDTO(img)
	return response.ResponseWithSuccessMessage(ctx, fiber.StatusCreated, "Updated File.", resp, nil)
}

func (c *StorageController) GetFileFromParameterHandler(ctx *fiber.Ctx) error {

	diskname := ctx.Params("disk")
	id := ctx.Params("id")

	img, err := c.service_storage.GetImage(diskname, id)
	if err != nil {
		return err
	}

	return ctx.SendStream(img.File, int(img.Size))
}
