package utils

import (
	"crypto/sha512"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/pkg/errors"

	"github.com/Chamistery/TestTask/internal/auth/model"
)

func GenerateToken(uuid, ip string, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Ip:   ip,
		Uuid: uuid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}

func HashToken(token string) (string, error) {
	shaHash := sha512.Sum512([]byte(token))
	hashedToken, err := bcrypt.GenerateFromPassword(shaHash[:], bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка хеширования токена: %v", err)
	}
	if len(hashedToken) == 0 {
		return "", fmt.Errorf("хеш не был сгенерирован")
	}
	return string(hashedToken), nil
}

func VerifyHashedToken(hashedToken string, candidateToken string) bool {
	candidateShaToken := sha512.Sum512([]byte(candidateToken))
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), candidateShaToken[:])
	return err == nil
}
