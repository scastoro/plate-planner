package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (apiCfg *apiConfig) ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
		headerArray := strings.Split(authorizationHeader, " ")
		if len(headerArray) < 2 {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
		headerToken := headerArray[1]

		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			respondWithError(w, http.StatusInternalServerError, "Server error")
			return
		}

		token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("could not validate token")
			}
			return []byte(secret), nil
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user, ok := claims["user"].(UserModel)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
				return
			}

			ctx := context.WithValue(req.Context(), UserKey, user)
			req = req.WithContext(ctx)
			next(w, req)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Not authorized to view this resource")
			return
		}
	})
}
