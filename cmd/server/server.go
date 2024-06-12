package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kouzoh/kokukuma-fido/internal/server"
)

func main() {
	srv := server.NewServer()

	r := mux.NewRouter()
	r.Use(handlers.CORS(
		handlers.AllowedMethods([]string{"POST", "GET"}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	))

	r.HandleFunc("/getRequest", srv.GetRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/verifyResponse", srv.VerifyResponse).Methods("POST", "OPTIONS")

	serverAddress := ":8080"
	log.Println("starting fido server at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
