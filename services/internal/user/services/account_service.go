package services

import (
	"services/internal/user/models"
	"services/internal/user/repositories"

	"golang.org/x/crypto/bcrypt"
)

// Account servisi yapısı tanımıdır.
// It is a Account Service struct.
type AccountService struct {
	repository_account repositories.AccountRepository
}

// Yeni Account Servisi yaratır. Hata varsa sonuç olarak yansıtabilir.
// You can create Account Service. If there are any errors, can return its.
func NewAccountService(repository repositories.AccountRepository) (*AccountService, error) {

	// Eğer eklenecek başka altyapılar varsa buraya eklenebilir.
	// If you want to add another infrastructure, you can add here.
	return &AccountService{
		repository_account: repository,
	}, nil
}

func (a *AccountService) ConnectWithUser(user *models.User, username string, password string) (*models.Account, error) {

	// Hash Password
	password, err := a.generateHashWithString(password, 14)
	if err != nil {
		return nil, err
	}

	accountData := models.NewAccount(username, password, user)
	account, err := a.repository_account.Create(accountData)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountService) GetAccountWithUsername(username string) (*models.Account, error) {
	account, err := a.repository_account.FindOneWithParameters("username", username)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountService) generateHashWithString(str string, length int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), length)
	if err != nil {
		return "", err
	}
	password := string(hash)

	return password, nil
}

func (a *AccountService) CompareHashAndPassword(hpass string, pass string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}
