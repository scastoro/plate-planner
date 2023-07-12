package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/scastoro/plate-planner-api/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port env variable not found!")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Could not get DB Url from environment variable!")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Post("/log-in", apiCfg.handlerLoginUser)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.ValidateTokenMiddleware(apiCfg.handlerGetUserById))
	v1Router.Get("/users/{userId}/workouts", apiCfg.ValidateTokenMiddleware(apiCfg.handlerGetWorkoutsByUserId))
	v1Router.Put("/users/{userId}/workouts", apiCfg.ValidateTokenMiddleware(apiCfg.handlerUpdateWorkout))

	v1Router.Post("/workouts", apiCfg.ValidateTokenMiddleware(apiCfg.handlerCreateWorkout))
	v1Router.Get("/workouts", apiCfg.ValidateTokenMiddleware(apiCfg.handlerGetWorkoutById))
	v1Router.Get("/workouts/{workoutId}/sets", apiCfg.ValidateTokenMiddleware(apiCfg.handlerGetSetsByWorkoutId))
	v1Router.Put("/workouts/{workoutId}/sets", apiCfg.ValidateTokenMiddleware(apiCfg.handlerUpdateSet))
	v1Router.Get("/workouts-with-sets/{workoutId}", apiCfg.ValidateTokenMiddleware(apiCfg.handlerGetWorkoutsByIdWithSets))

	v1Router.Post("/sets", apiCfg.ValidateTokenMiddleware(apiCfg.handlerCreateSet))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server running on port: %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
