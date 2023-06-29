package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	dat, err := json.Marshal(struct{}{})
	if err != nil {
		log.Println("Failed to marshal json")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
