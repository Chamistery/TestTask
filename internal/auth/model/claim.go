package model

const (
	CreatePath  = "/Create"
	RefreshPath = "/Refresh"
)

type UserClaims struct {
	jwt.StandardClaims
	Ip string `json:"ip"`
}
