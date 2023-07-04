package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func (apiCfg *apiConfig) ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			respondWithError(w, http.StatusInternalServerError, "Server error")
			return
		}

		token, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("could not validate token")
			}
			return []byte(secret), nil
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
		if token.Valid {
			next(w, req)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
	})
}
