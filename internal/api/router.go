package api

import (
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/api/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	healthPath       = "/health"
	versionPath      = "/version"
	housesPath       = "/houses"
	extrasPath       = "/extras"
	reservationPath  = "/reservation"
	verificationPath = "/verification"
	bathhousesPath   = "/bathhouses"
	idPath           = "/{id}"
	emptyPath        = ""
)

type Middlewares struct {
	PanicRecovery mux.MiddlewareFunc
}

type IReservations interface {
	CreateReservation(w http.ResponseWriter, r *http.Request)
	GetAvailableHouses(w http.ResponseWriter, r *http.Request)
}

type IHouses interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IBathhouses interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByHouse(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IExtras interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IVerification interface {
	VerifyIdentity(w http.ResponseWriter, r *http.Request)
}

type IGeneral interface {
	Health(w http.ResponseWriter, r *http.Request)
	Version(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	Reservations IReservations
	Houses       IHouses
	Bathhouses   IBathhouses
	Extras       IExtras
	Verification IVerification
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

	r.HandleFunc(healthPath, dep.Handlers.General.Health)
	r.HandleFunc(versionPath, dep.Handlers.General.Version)

	r.HandleFunc(verificationPath, dep.Handlers.Verification.VerifyIdentity).Methods("POST")

	reservations := r.PathPrefix(reservationPath).Subrouter()
	reservations.HandleFunc(emptyPath, dep.Handlers.Reservations.GetAvailableHouses).Methods(http.MethodGet)
	reservations.HandleFunc(emptyPath, dep.Handlers.Reservations.CreateReservation).Methods(http.MethodPost)

	houses := r.PathPrefix(housesPath).Subrouter()
	houses.HandleFunc(emptyPath, dep.Handlers.Houses.Add).Methods(http.MethodPost)
	houses.HandleFunc(idPath, dep.Handlers.Houses.Update).Methods(http.MethodPut)
	houses.HandleFunc(idPath, dep.Handlers.Houses.Delete).Methods(http.MethodDelete)
	houses.HandleFunc(emptyPath, dep.Handlers.Houses.GetAll).Methods(http.MethodGet)

	bathhouses := r.PathPrefix(bathhousesPath).Subrouter()
	bathhouses.HandleFunc(emptyPath, dep.Handlers.Bathhouses.Add).Methods(http.MethodPost)
	bathhouses.HandleFunc(idPath, dep.Handlers.Bathhouses.Update).Methods(http.MethodPut)
	bathhouses.HandleFunc(idPath, dep.Handlers.Bathhouses.Delete).Methods(http.MethodDelete)
	bathhouses.HandleFunc(emptyPath, dep.Handlers.Bathhouses.GetAll).Methods(http.MethodGet)
	bathhouses.HandleFunc(idPath, dep.Handlers.Bathhouses.GetByHouse).Methods(http.MethodGet)

	extras := r.PathPrefix(extrasPath).Subrouter()
	extras.HandleFunc(emptyPath, dep.Handlers.Extras.Add).Methods(http.MethodPost)
	extras.HandleFunc(idPath, dep.Handlers.Extras.Update).Methods(http.MethodPut)
	extras.HandleFunc(idPath, dep.Handlers.Extras.Delete).Methods(http.MethodDelete)
	extras.HandleFunc(emptyPath, dep.Handlers.Extras.GetAll).Methods(http.MethodGet)

	return middleware.WithCORS(r)
}
