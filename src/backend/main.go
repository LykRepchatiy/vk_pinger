package main

import (
	database "backend/database"
	handlers "backend/handlers"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	env := database.ParseEnv()
	mux := http.NewServeMux()
	mux.HandleFunc("/putStatus", handlers.PutStatus)
	mux.HandleFunc("/containerList", handlers.ContainerList)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":"+env.Port, handler)
}
