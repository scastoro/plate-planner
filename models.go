package main

import (
	"fmt"
	"time"

	"github.com/scastoro/plate-planner-api/internal/database"
)

type UserModel struct {
	ID         int32
	Name       string
	BodyWeight string
	Username   string
	Email      string
}

func convertDbUserToUser(dbUser database.AdminUser) UserModel {
	return UserModel{
		ID:         dbUser.ID,
		Name:       fmt.Sprintf("%v %v", dbUser.FirstName, dbUser.LastName),
		BodyWeight: dbUser.BodyWeight,
		Username:   dbUser.Username,
		Email:      dbUser.Email,
	}
}

func convertDbUsersToUsers(dbUsers []database.AdminUser) []UserModel {
	users := []UserModel{}
	for _, dbUser := range dbUsers {
		users = append(users, UserModel{
			ID:         dbUser.ID,
			Name:       fmt.Sprintf("%v %v", dbUser.FirstName, dbUser.LastName),
			BodyWeight: dbUser.BodyWeight,
			Username:   dbUser.Username,
			Email:      dbUser.Email,
		})
	}

	return users
}

type WorkoutModel struct {
	ID            int32
	StartTime     time.Time
	Duration      string
	TotalWeight   string
	TotalCalories int32
}

func convertDbWorkoutToWorkout(dbWorkout database.Workout) WorkoutModel {
	return WorkoutModel{
		ID:            dbWorkout.ID,
		StartTime:     dbWorkout.StartTime,
		Duration:      dbWorkout.Duration,
		TotalWeight:   dbWorkout.TotalWeight,
		TotalCalories: dbWorkout.TotalCalories,
	}
}

func convertDbWorkoutsToWorkouts(dbWorkouts []database.GetWorkoutsByUserIdDescRow) []WorkoutModel {
	workouts := []WorkoutModel{}
	for _, dbWorkout := range dbWorkouts {
		workouts = append(workouts, WorkoutModel{
			ID:            dbWorkout.ID,
			StartTime:     dbWorkout.StartTime,
			Duration:      dbWorkout.Duration,
			TotalWeight:   dbWorkout.TotalWeight,
			TotalCalories: dbWorkout.TotalCalories,
		})
	}

	return workouts
}

type SetModel struct {
	ID        int32
	Exercise  string
	Count     int32
	Intensity database.Intensity
	Type      string
	Weight    string
}

func convertDbSetToSet(dbSet database.Set) SetModel {
	return SetModel{
		ID:        dbSet.ID,
		Exercise:  dbSet.Exercise,
		Count:     dbSet.Count,
		Intensity: dbSet.Intensity,
		Type:      dbSet.Type,
		Weight:    dbSet.Weight,
	}
}

func convertDbSetsToSets(dbSets []database.Set) []SetModel {
	sets := []SetModel{}
	for _, dbSet := range dbSets {
		sets = append(sets, SetModel{
			ID:        dbSet.ID,
			Exercise:  dbSet.Exercise,
			Count:     dbSet.Count,
			Intensity: dbSet.Intensity,
			Type:      dbSet.Type,
			Weight:    dbSet.Weight,
		})
	}

	return sets
}

type envelope map[string]interface{}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	LastPage     int
	TotalRecords int
}
