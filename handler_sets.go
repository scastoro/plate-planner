package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/scastoro/plate-planner-api/internal/database"
)

func (apiCfg *apiConfig) handlerCreateSet(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		WorkoutId int32  `json:"workout_id"`
		Count     int32  `json:"count"`
		Intensity string `json:"intensity"`
		Type      string `json:"type"`
		Weight    string `json:"weight"`
		Exercise  string `json:"exercise"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	set, err := apiCfg.DB.CreateSet(r.Context(), database.CreateSetParams{
		Exercise:  params.Exercise,
		Count:     params.Count,
		Intensity: database.Intensity(params.Intensity),
		Type:      params.Type,
		Weight:    params.Weight,
		WorkoutID: params.WorkoutId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating set: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, convertDbSetToSet(set))
}

func (apiCfg *apiConfig) handlerGetSetsByWorkoutId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("workout_id")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	workoutId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting workout id param to num")
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

	sets, err := apiCfg.DB.GetSetsByWorkoutIdDesc(r.Context(), database.GetSetsByWorkoutIdDescParams{
		WorkoutID:  int32(workoutId),
		OrderByCol: "created_at",
		Offset:     int32(parsedOffset),
		Limit:      int32(parsedLimit),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbSetsToSets(sets))
}