package model

import "github.com/dgrijalva/jwt-go"

const (
	CreatePath  = "/Create"
	RefreshPath = "/Refresh"
)

type UserClaims struct {
	jwt.StandardClaims
	Ip   string `json:"ip"`
	Uuid string `json:"uuid"`
}
