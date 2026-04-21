package dto

type CreateUserAuthDataDTOInput struct {
	UserID         int64
	HashedPassword string
}

type CreateUserAuthDataDTOOutput struct {
	ID             int64
	UserID         int64
	HashedPassword string
}

type UpdateUserAuthDataDTOInput struct {
	HashedPassword string
}

type GetUserAuthDataByUserIDDTOOutput struct {
	ID             int64
	HashedPassword string
}

type SetUserAuthDataByUserIDDTOInput struct {
	ID             *int64
	UserID         int64
	HashedPassword *string
}
