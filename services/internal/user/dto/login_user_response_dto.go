package dto

type LoginUserResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewLoginUserResponseDTO(accessToken string, refreshToken string) *LoginUserResponseDTO {
	return &LoginUserResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
