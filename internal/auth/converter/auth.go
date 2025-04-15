package converter

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(refresh_token string) string {
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refresh_token), bcrypt.DefaultCost)
	return string(hashedToken)
}
