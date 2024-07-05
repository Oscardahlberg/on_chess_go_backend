package main

import (
	"log"
	"net/http"

	"gobackend/database"
	"gobackend/routes"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Println("Error initializing db: ", err)
	} else {
		log.Println("MongoDB link started")
	}

	routing := routes.SetupRoutes()

	log.Println("Server started on https://localhost:3001/")
	if err := http.ListenAndServe(":3001", routing); err != nil {
		log.Fatalf("error: %v", err)
	}
}
