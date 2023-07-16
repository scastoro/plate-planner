package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/scastoro/plate-planner-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
		Email     string `json:"email"`
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
		Email:        params.Email,
		Password:     hash,
		Lastloggedin: time.Now(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error saving user to the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbUserToUser(user))
}

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserWithPermissionsByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Username or password incorrect")
		return
	}

	ok := CheckPassword(params.Password, user.Password)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Username or password incorrect")
		return
	}

	userModel := convertDbUserWithPermsToUserWithPerms(user)
	token, err := CreateToken(userModel)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	respondWithJson(w, http.StatusOK, envelope{"token": token})
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

	respondWithJson(w, 200, convertDbUserToUser(user))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (apiCfg apiConfig) handlerGetUserByIdWithPerms(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting user id param to num")
		return
	}

	user, err := apiCfg.DB.GetUserWithPermissions(r.Context(), int32(userId))
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, convertDbUserWithPermsToUserWithPerms(user))
}
