package model

type CreateModel struct {
	Uuid         string
	RefreshToken string
	Guid         string
}

type RefreshModel struct {
	Uuid            string
	RefreshTokenOld string
	RefreshTokenNew string
}
