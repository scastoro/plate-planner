package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	id := r.URL.Query().Get("user_id")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting user id param to num")
		return
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		respondWithError(w, 400, "Error getting the offset from query string")
		return
	}
	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, 500, "Error converting offset param to num")
		return
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		respondWithError(w, 400, "Error getting the limit from query string")
		return
	}
	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, 500, "Error converting limit param to num")
		return
	}

	workouts, err := apiCfg.DB.GetWorkoutsByUserIdDesc(r.Context(), database.GetWorkoutsByUserIdDescParams{
		UserID:     int32(userId),
		OrderByCol: "StartTime",
		Offset:     int32(parsedOffset),
		Limit:      int32(parsedLimit),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbWorkoutsToWorkouts(workouts))
}
