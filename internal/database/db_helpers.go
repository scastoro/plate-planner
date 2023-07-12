package database

import "context"

type WorkoutWithSets struct {
	Workout
	Sets         []Set
	TotalRecords int64
}

func (q *Queries) GetWorkoutsSetsHelper(ctx context.Context, params GetWorkoutsByIdWithSetsParams) (WorkoutWithSets, error) {

	var workoutsWithSets = WorkoutWithSets{}

	workouts, err := q.GetWorkoutsByIdWithSets(ctx, params)
	if err != nil {
		return workoutsWithSets, err
	}
	if len(workouts) < 1 {
		return workoutsWithSets, nil
	}

	firstRow := workouts[0]
	workoutsWithSets.Workout = Workout{
		ID:            firstRow.ID,
		StartTime:     firstRow.StartTime,
		Duration:      firstRow.Duration,
		TotalWeight:   firstRow.TotalWeight,
		TotalCalories: firstRow.TotalCalories,
		UserID:        firstRow.UserID,
	}
	workoutsWithSets.TotalRecords = firstRow.Count

	sets := []Set{}
	for _, row := range workouts {
		sets = append(sets, Set{
			ID:        row.ID_2,
			Exercise:  row.Exercise,
			Count:     row.Count_2,
			Intensity: row.Intensity,
			Type:      row.Type,
			Weight:    row.Weight,
			WorkoutID: row.WorkoutID,
		})
	}

	workoutsWithSets.Sets = sets

	return workoutsWithSets, nil

}
