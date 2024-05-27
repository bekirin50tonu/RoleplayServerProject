package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage struct {
	Filename  string `bson:"filename"`
	Disk      string `bson:"disk"`
	Size      int64  `bson:"size"`
	Extension string `bson:"extension"`

	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

func NewStorage(filename, disk, extension string, size int64) *Storage {

	return &Storage{
		Filename:  filename,
		Disk:      disk,
		Size:      size,
		Extension: extension,
		ID:        primitive.NewObjectID(),
	}

}
