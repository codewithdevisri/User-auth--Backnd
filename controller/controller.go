package controller

import (
	authorization "github.com/Ayan25844/netflix/middleware"
	"github.com/Ayan25844/netflix/service"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	accountsubrouter := router.PathPrefix("/api").Subrouter()
	accountsubrouter.HandleFunc(`/user/login`, service.Login).Methods(`POST`)
	accountsubrouter.HandleFunc("/user/signup", service.CreateUser).Methods("POST")

	// Group for endpoints accessible only by Admins
	adminSubrouter := router.PathPrefix("/api").Subrouter()
	adminSubrouter.Use(authorization.ValidateToken)
	adminSubrouter.Use(authorization.Authorization([]string{"ADMIN"}))

	adminSubrouter.HandleFunc("/users", service.GetAll).Methods("GET")
	adminSubrouter.HandleFunc("/user", service.CreateUser).Methods("POST")
	adminSubrouter.HandleFunc("/user/{id}", service.UpdateUser).Methods("PUT")
	adminSubrouter.HandleFunc("/user/{id}", service.DeleteUser).Methods("DELETE")
	adminSubrouter.HandleFunc("/deletealluser", service.DeleteAll).Methods("DELETE")

	// Group for endpoints accessible by USER and ADMIN both
	userSubrouter := router.PathPrefix("/api").Subrouter()
	userSubrouter.Use(authorization.ValidateToken)
	userSubrouter.Use(authorization.Authorization([]string{"USER", "ADMIN"}))

	userSubrouter.HandleFunc("/user/{id}", service.FindById).Methods("GET")

	return router
}
