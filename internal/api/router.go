package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Middlewares struct {
	PanicRecovery mux.MiddlewareFunc
}

type IReservations interface {
	BookAHouse(w http.ResponseWriter, r *http.Request)
}

type IHouses interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IGeneral interface {
	Health(w http.ResponseWriter, r *http.Request)
	Version(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	Reservations IReservations
	Houses       IHouses
	General      IGeneral
}

type RouterDependencies struct {
	Handlers    Handlers
	Middlewares Middlewares
}

func NewRouter(dep RouterDependencies) http.Handler {
	r := mux.NewRouter()

	r.Use(dep.Middlewares.PanicRecovery.Middleware)

	r.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
	})

	r.HandleFunc("/health", dep.Handlers.General.Health)
	r.HandleFunc("/version", dep.Handlers.General.Version)

	reservations := r.PathPrefix("/reservations").Subrouter()
	reservations.HandleFunc("/house", dep.Handlers.Reservations.BookAHouse).Methods(http.MethodPost)

	houses := r.PathPrefix("/houses").Subrouter()
	houses.HandleFunc("", dep.Handlers.Houses.Add).Methods(http.MethodPost)
	houses.HandleFunc("/{id}", dep.Handlers.Houses.Update).Methods(http.MethodPut)
	houses.HandleFunc("/{id}", dep.Handlers.Houses.Delete).Methods(http.MethodDelete)
	houses.HandleFunc("", dep.Handlers.Houses.GetAll).Methods(http.MethodGet)

	return r
}
