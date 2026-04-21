package dto

import "time"

type CreateUserDTOInput struct {
	Username string
}

type CreateUserDTOOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}

type UpdateUserDTOInput struct {
	Username          *string
	ProfilePictureURL *string
}

type GetUserByIDDTOOutput struct {
	Username          string
	CreatedAt         time.Time
	ProfilePictureURL string
}

type SetUserDTOInput struct {
	ID                int64
	Username          *string
	CreatedAt         *time.Time
	ProfilePictureURL *string
}
