package dto

type LoginDTOInput struct {
	Username string
	Password string
}

type LoginDTOOutput struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
	CacheErrors  struct {
		RefreshToken error
	}
}

type RegisterDTOInput struct {
	Username string
	Password string
}

type RegisterDTOOutput struct {
	AccessToken  string
	RefreshToken string
	CacheErrors  struct {
		User         error
		UserAuthData error
		RefreshToken error
	}
}

type RefreshDTOInput struct {
	RefreshToken string
}

type RefreshDTOOutput struct {
	AccessToken  string
	RefreshToken string
}

type LogoutDTOInput struct {
	RefreshToken string
}
