package dto

type RegisterDTOInput struct {
	Username string
	Password string
}

type RegisterDTOOutput struct {
	AccessToken  string
	RefreshToken string
}
