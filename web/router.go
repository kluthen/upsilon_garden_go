package web

import (
	"log"
	"net/http"
	"upsilon_garden_go/config"
	"upsilon_garden_go/web/garden_controller"
	"upsilon_garden_go/web/plant_controller"

	"github.com/gorilla/mux"
)

// RouterSetup Prepare routing.
func RouterSetup() *mux.Router {
	r := mux.NewRouter()

	// CRUD /gardens
	r.HandleFunc("/gardens", garden_controller.Index).Methods("GET")
	r.HandleFunc("/gardens/{id}", garden_controller.Show).Methods("GET")
	r.HandleFunc("/gardens", garden_controller.Create).Methods("POST")
	r.HandleFunc("/gardens/{id}", garden_controller.Update).Methods("PUT")
	r.HandleFunc("/gardens/{id}", garden_controller.Delete).Methods("DELETE")

	// Hydro functions
	r.HandleFunc("/gardens/{id}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	r.HandleFunc("/gardens/{id}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")

	// CRUD /gardens/:id/plant
	plantRouter := r.PathPrefix("/gardens/{gid}/plants").Subrouter()
	plantRouter.HandleFunc("/", plant_controller.Index).Methods("GET")
	plantRouter.HandleFunc("/{pid}", plant_controller.Show).Methods("GET")
	plantRouter.HandleFunc("/", plant_controller.Create).Methods("POST")
	plantRouter.HandleFunc("/{pid}", plant_controller.Update).Methods("PUT")
	plantRouter.HandleFunc("/{pid}", plant_controller.Delete).Methods("DELETE")

	// JSON Access ...

	jsonAPI := r.PathPrefix("/api").Subrouter()

	// CRUD /gardens
	jsonAPI.HandleFunc("/gardens", garden_controller.Index).Methods("GET")
	jsonAPI.HandleFunc("/gardens/{id}", garden_controller.Show).Methods("GET")
	jsonAPI.HandleFunc("/gardens", garden_controller.Create).Methods("POST")
	jsonAPI.HandleFunc("/gardens/{id}", garden_controller.Update).Methods("PUT")
	jsonAPI.HandleFunc("/gardens/{id}", garden_controller.Delete).Methods("DELETE")

	// Hydro functions
	jsonAPI.HandleFunc("/gardens/{id}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	jsonAPI.HandleFunc("/gardens/{id}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")

	// CRUD /gardens/:id/plant
	APIPlantRouter := jsonAPI.PathPrefix("/gardens/{gid}/plants").Subrouter()
	APIPlantRouter.HandleFunc("/", plant_controller.Index).Methods("GET")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Show).Methods("GET")
	APIPlantRouter.HandleFunc("/", plant_controller.Create).Methods("POST")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Update).Methods("PUT")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Delete).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.STATIC_FILES))))
	r.Use(loggingMw)
	return r
}

func loggingMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Web: Received request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// ListenAndServe start listing http server
func ListenAndServe(router *mux.Router) {
	log.Printf("Web: Started server on 127.0.0.1:80 and listening ... ")
	http.ListenAndServe(":80", router)
}
