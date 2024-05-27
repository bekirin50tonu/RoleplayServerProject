package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Name      string             `bson:"name"`
	Lastname  string             `bson:"lastname"`
	Email     string             `bson:"email"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

func NewUser(name string, lastname string, email string) *User {
	return &User{
		Name:     name,
		Lastname: lastname,
		Email:    email,
		ID:       primitive.NewObjectID(),
	}
}
