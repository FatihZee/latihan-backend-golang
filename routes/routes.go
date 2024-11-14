package routes

import (
	"database/sql"
	"myapp/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Routes untuk registrasi dan login
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.Register(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.Login(db, w, r)
	}).Methods("POST")

	// Routes untuk user
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUsers(db, w, r)
	}).Methods("GET")
	r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateUser(db, w, r)
	}).Methods("PUT")
	r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteUser(db, w, r)
	}).Methods("DELETE")

	return r
}
