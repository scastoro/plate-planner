package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
	id := chi.URLParam(r, "workoutId")
	if id == "" {
		respondWithError(w, 400, "Error getting the user id from query string")
		return
	}
	workoutId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, "Error converting workout id param to num")
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

	sets, err := apiCfg.DB.GetSetsByWorkoutIdDesc(r.Context(), database.GetSetsByWorkoutIdDescParams{
		WorkoutID:  int32(workoutId),
		OrderByCol: "created_at",
		Offset:     int32(offset),
		Limit:      int32(parsedCount),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting the user from the database: %v", err))
		return
	}

	respondWithJson(w, 200, convertDbSetsToSets(sets))
}

func (apiCfg *apiConfig) handlerUpdateSet(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Exercise  string `json:"exercise"`
		Count     int32  `json:"count"`
		Intensity string `json:"intensity"`
		Type      string `json:"type"`
		Weight    string `json:"weight"`
		ID        int32  `json:"set_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing json body")
		return
	}

	set, err := apiCfg.DB.UpdateSetById(r.Context(), database.UpdateSetByIdParams{
		Exercise:  params.Exercise,
		Count:     params.Count,
		Intensity: database.Intensity(params.Intensity),
		Type:      params.Type,
		Weight:    params.Weight,
		ID:        params.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating workout: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, convertDbSetToSet(set))
}
