package model

type CreateRequest struct {
	Guid string `json:"guid"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshRequest) GetRefreshToken() string {
	return r.RefreshToken
}
