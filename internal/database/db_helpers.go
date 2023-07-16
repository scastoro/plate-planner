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

type UserWithPermissions struct {
	AdminUser
	Role        string
	Permissions []Permission
}

func (q *Queries) GetUserWithPermissions(ctx context.Context, id int32) (UserWithPermissions, error) {
	userWithPermissions := UserWithPermissions{}

	users, err := q.GetUserByIdWithPerms(ctx, id)
	if err != nil {
		return UserWithPermissions{}, err
	}
	if len(users) < 1 {
		return UserWithPermissions{}, nil
	}

	firstRow := users[0]

	userWithPermissions.Role = firstRow.Role
	userWithPermissions.AdminUser = AdminUser{
		FirstName:    firstRow.FirstName,
		LastName:     firstRow.LastName,
		BodyWeight:   firstRow.BodyWeight,
		ID:           firstRow.ID,
		RoleID:       firstRow.RoleID,
		Lastloggedin: firstRow.Lastloggedin,
		Email:        firstRow.Email,
		Password:     firstRow.Password,
		Username:     firstRow.Username,
	}

	permissions := []Permission{}

	for _, user := range users {
		permissions = append(permissions, Permission{
			Action:   user.Permission,
			Resource: user.Resource,
		})
	}

	userWithPermissions.Permissions = permissions

	return userWithPermissions, nil
}

func (q *Queries) GetUserWithPermissionsByEmail(ctx context.Context, email string) (UserWithPermissions, error) {
	userWithPermissions := UserWithPermissions{}

	users, err := q.GetUserByEmailWithPerms(ctx, email)
	if err != nil {
		return UserWithPermissions{}, err
	}
	if len(users) < 1 {
		return UserWithPermissions{}, nil
	}

	firstRow := users[0]

	userWithPermissions.Role = firstRow.Role
	userWithPermissions.AdminUser = AdminUser{
		FirstName:    firstRow.FirstName,
		LastName:     firstRow.LastName,
		BodyWeight:   firstRow.BodyWeight,
		ID:           firstRow.ID,
		RoleID:       firstRow.RoleID,
		Lastloggedin: firstRow.Lastloggedin,
		Email:        firstRow.Email,
		Password:     firstRow.Password,
		Username:     firstRow.Username,
	}

	permissions := []Permission{}

	for _, user := range users {
		permissions = append(permissions, Permission{
			Action:   user.Permission,
			Resource: user.Resource,
		})
	}

	userWithPermissions.Permissions = permissions

	return userWithPermissions, nil
}
