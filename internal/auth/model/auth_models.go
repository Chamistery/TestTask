package model

type CreateModel struct {
	RefreshToken string
	Guid         string
}

type RefreshModel struct {
	RefreshTokenOld string
	RefreshTokenNew string
}
