package services

import (
	"log"
	config "services/internal/user/config/jwt"
	"services/internal/user/dto"
	"services/internal/user/models"
	"services/internal/user/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SessionService struct {
	repository_session *repositories.SessionRepository
}

func NewSessionService(repo *repositories.SessionRepository) (*SessionService, error) {

	return &SessionService{
		repository_session: repo,
	}, nil
}

func (s *SessionService) NewAccessToken(account *models.Account) (string, error) {
	accessDTO := dto.NewJWTTokenDTO(account.Username, account.ID, time.Now().Add(config.ACCESSTOKENTIME))
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessDTO)
	return accessToken.SignedString([]byte(config.ACCESSTOKENSECRET))

}
func (s *SessionService) NewRefreshToken(account *models.Account) (string, error) {
	refreshDTO := dto.NewJWTTokenDTO(account.Username, account.ID, time.Now().Add(config.ACCESSTOKENTIME))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshDTO)
	return refreshToken.SignedString([]byte(config.REFRESHTOKENSECRET))
}

func (s *SessionService) ParseAccessToken(accessToken string) (*dto.JWTTokenDTO, error) {
	parsed, err := jwt.ParseWithClaims(accessToken, &dto.JWTTokenDTO{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.ACCESSTOKENSECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return parsed.Claims.(*dto.JWTTokenDTO), nil
}

func (s *SessionService) ParseRefreshToken(refreshToken string) (*dto.JWTTokenDTO, error) {
	parsed, err := jwt.ParseWithClaims(refreshToken, &dto.JWTTokenDTO{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.REFRESHTOKENSECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return parsed.Claims.(*dto.JWTTokenDTO), nil
}

func (s *SessionService) CreateSession(account *models.Account) (*models.Session, error) {

	accessToken, err := s.NewAccessToken(account)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.NewRefreshToken(account)
	if err != nil {
		return nil, err
	}

	session := models.NewSession(account, accessToken, refreshToken)

	session, err = s.repository_session.UpdateOrCreateSession(session, "account")
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetAccountWithAccessToken(accessToken string) (*models.Account, error) {

	session, err := s.repository_session.FindOneWithParameter("access_token", accessToken)
	if err != nil {
		log.Default().Print(err, accessToken)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Access Token Expired.", err.Error())
	}

	_, err = s.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	return session.Account, nil
}

func (s *SessionService) GetAccountWithTokens(refreshToken string) (*models.Account, error) {

	session, err := s.repository_session.FindOneWithParameter("refresh_token", refreshToken)
	if err != nil {
		log.Default().Print(err, refreshToken)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Refresh Token Expired.", err.Error())
	}

	_, err = s.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return session.Account, nil
}
