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

type IGeneral interface {
	Health(w http.ResponseWriter, r *http.Request)
	Version(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	Reservations IReservations
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

	Internal := r.PathPrefix("/internal").Subrouter()

	weather := Internal.PathPrefix("/reservations").Subrouter()
	weather.HandleFunc("/house", dep.Handlers.Reservations.BookAHouse).Methods(http.MethodPost)

	return r
}
