package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/scastoro/plate-planner-api/internal/database"
)

func (apiCfg *apiConfig) handlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserId string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	userId, err := strconv.Atoi(params.UserId)
	if err != nil {
		respondWithError(w, 500, "Error converting param to num")
		return
	}

	workout, err := apiCfg.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		StartTime:     time.Now(),
		Duration:      "0.0",
		TotalWeight:   "0.0",
		TotalCalories: 0,
		UserID:        int32(userId),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error saving workout to the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbWorkoutToWorkout(workout))
}

func (apiCfg *apiConfig) handlerGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("workout_id")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	workoutId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting param to num")
		return
	}

	workout, err := apiCfg.DB.GetWorkoutById(r.Context(), int32(workoutId))
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbWorkoutToWorkout(workout))
}

func (apiCfg *apiConfig) handlerGetWorkoutsByUserId(w http.ResponseWriter, r *http.Request) {
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
	count := r.URL.Query().Get("count")
	if count == "" {
		respondWithError(w, 400, "Error getting the count from query string")
		return
	}
	parsedCount, err := strconv.Atoi(count)
	if err != nil {
		respondWithError(w, 500, "Error converting count param to num")
		return
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		respondWithError(w, 400, "Error getting the page from query string")
		return
	}
	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		respondWithError(w, 500, "Error converting page param to num")
		return
	}

	offset := parsedCount * (parsedPage - 1)

	workouts, err := apiCfg.DB.GetWorkoutsByUserIdDesc(r.Context(), database.GetWorkoutsByUserIdDescParams{
		UserID:     int32(userId),
		OrderByCol: "StartTime",
		Offset:     int32(offset),
		Limit:      int32(parsedCount),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	pageData := Metadata{
		PageSize:    parsedCount,
		CurrentPage: parsedPage,
		FirstPage:   1,
		LastPage:    1,
	}
	if len(workouts) > 0 {
		totalRecords := workouts[0].Count
		pageData.TotalRecords = int(totalRecords)
		pageData.LastPage = int(math.Ceil(float64(totalRecords) / float64(parsedCount)))
	}

	respondWithJson(w, 200, envelope{"metadata": pageData, "records": convertDbWorkoutsToWorkouts(workouts)})
}

func (apiCfg *apiConfig) handlerUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Duration      string `json:"duration"`
		TotalWeight   string `json:"total_weight"`
		TotalCalories int32  `json:"total_calories"`
		Id            int32  `json:"workout_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing json body")
		return
	}

	workout, err := apiCfg.DB.UpdateWorkoutById(r.Context(), database.UpdateWorkoutByIdParams{
		ID:            params.Id,
		Duration:      params.Duration,
		TotalWeight:   params.TotalWeight,
		TotalCalories: params.TotalCalories,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating workout: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, convertDbWorkoutToWorkout(workout))
}

func (apiCfg *apiConfig) handlerGetWorkoutsByIdWithSets(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "workoutId")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	workoutId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting user id param to num")
		return
	}
	count := r.URL.Query().Get("count")
	if count == "" {
		respondWithError(w, 400, "Error getting the count from query string")
		return
	}
	parsedCount, err := strconv.Atoi(count)
	if err != nil {
		respondWithError(w, 500, "Error converting count param to num")
		return
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		respondWithError(w, 400, "Error getting the page from query string")
		return
	}
	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		respondWithError(w, 500, "Error converting page param to num")
		return
	}
	offset := parsedCount * (parsedPage - 1)

	workoutWithSets, err := apiCfg.DB.GetWorkoutsSetsHelper(r.Context(), database.GetWorkoutsByIdWithSetsParams{
		ID:     int32(workoutId),
		Limit:  int32(parsedCount),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting workouts with sets: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, workoutWithSets)
}
