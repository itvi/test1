package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte("hxahxhfqekjkgloafafk.dfakf")

func GenerateToken() (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	})
	token, err := claims.SignedString(key)
	if err != nil {
		fmt.Printf("Signed error: %s", err.Error())
		return "", err
	}
	return token, nil
}
