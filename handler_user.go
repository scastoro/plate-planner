package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/scastoro/plate-planner-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	hash, err := HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error hashing password: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		BodyWeight:   "0.0",
		Username:     fmt.Sprintf("%v%v", params.FirstName, params.LastName),
		Email:        "test",
		Password:     hash,
		Lastloggedin: time.Now(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error saving user to the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	num, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting param to num")
		return
	}

	user, err := apiCfg.DB.GetUserById(r.Context(), int32(num))
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, 200, user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
