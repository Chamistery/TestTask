package model

import "github.com/Chamistery/TestTask/internal/auth/auth_http"

type CreateRequest struct {
	Guid string `json:"Guid"`
}

type Response struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

type RefreshRequest struct {
	refresh_token string `json:"refresh_token"`
}

type Implementation struct {
	authService auth_http.AuthService
}

func (r *RefreshRequest) GetRefreshToken() string {
	return r.refresh_token
}
