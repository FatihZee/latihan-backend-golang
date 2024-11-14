package main

import (
	"log"
	"myapp/config"
	"myapp/routes"
	"net/http"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	r := routes.SetupRoutes(db)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
