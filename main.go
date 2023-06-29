package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port env variable not found!")
	}

	router := chi.NewRouter()

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server running on port: %s", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
