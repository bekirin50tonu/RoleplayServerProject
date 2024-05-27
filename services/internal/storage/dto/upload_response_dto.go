package dto

import (
	"services/internal/storage/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadResponseDTO struct {
	ID        primitive.ObjectID `json:"id"`
	Filename  string             `json:"filename"`
	Disk      string             `json:"disk"`
	Size      int64              `json:"size"`
	Extension string             `json:"extension"`
}

func NewUploadResponseDTO(model *models.Storage) UploadResponseDTO {

	return UploadResponseDTO{
		ID:        model.ID,
		Filename:  model.Filename,
		Disk:      model.Disk,
		Size:      model.Size,
		Extension: model.Extension,
	}
}
