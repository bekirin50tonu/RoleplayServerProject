package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	User      *User              `bson:"user._id,omitempty"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

func NewAccount(username string, password string, user *User) *Account {

	return &Account{
		Username: username,
		Password: password,
		User:     user,
		ID:       primitive.NewObjectID(),
	}

}
