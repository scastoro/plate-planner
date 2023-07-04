package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user UserModel) (string, error) {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return "", errors.New("could not find secret key in environment")
	}
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user"] = user
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return token, nil
}
