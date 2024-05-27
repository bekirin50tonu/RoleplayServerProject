package services

import (
	"services/internal/user/models"
	"services/internal/user/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Account servisi yapısı tanımıdır.
// It is a Account Service struct.
type UserService struct {
	repository_user repositories.UserRepository
}

// Yeni Account Servisi yaratır. Hata varsa sonuç olarak yansıtabilir.
// You can create Account Service. If there are any errors, can return its.
func NewUserService(repository repositories.UserRepository) (*UserService, error) {

	// Eğer eklenecek başka altyapılar varsa buraya eklenebilir.
	// If you want to add another infrastructure, you can add here.
	return &UserService{
		repository_user: repository,
	}, nil
}

func (u *UserService) CreateUserWithData(name string, lastname string, email string) (*models.User, error) {
	user := models.NewUser(name, lastname, email)
	user, err := u.repository_user.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) DeleteUserWithID(id primitive.ObjectID) error {

	err := u.repository_user.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
