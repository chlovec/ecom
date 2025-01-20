package auth

import (
	"ecom/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	claims := jwt.MapClaims{
		"userID": strconv.Itoa(userId),
		"expiresAt": time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(secret)
}