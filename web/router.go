package web

import (
	"log"
	"net/http"
	"upsilon_garden_go/config"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/web/garden_controller"
	"upsilon_garden_go/web/plant_controller"
	"upsilon_garden_go/web/tools"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// RouterSetup Prepare routing.
func RouterSetup() *mux.Router {
	r := mux.NewRouter()

	// CRUD /gardens
	r.HandleFunc("/gardens", garden_controller.Index).Methods("GET")
	r.HandleFunc("/gardens/{gid}", garden_controller.Show).Methods("GET")
	r.HandleFunc("/gardens", garden_controller.Create).Methods("POST")
	r.HandleFunc("/gardens/{gid}", garden_controller.Update).Methods("PUT")
	r.HandleFunc("/gardens/{gid}", garden_controller.Delete).Methods("DELETE")

	// Hydro functions
	r.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	r.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")

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
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Show).Methods("GET")
	jsonAPI.HandleFunc("/gardens", garden_controller.Create).Methods("POST")
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Update).Methods("PUT")
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Delete).Methods("DELETE")

	// Hydro functions
	jsonAPI.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	jsonAPI.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")

	// CRUD /gardens/:id/plant
	APIPlantRouter := jsonAPI.PathPrefix("/gardens/{gid}/plants").Subrouter()
	APIPlantRouter.HandleFunc("/", plant_controller.Index).Methods("GET")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Show).Methods("GET")
	APIPlantRouter.HandleFunc("/", plant_controller.Create).Methods("POST")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Update).Methods("PUT")
	APIPlantRouter.HandleFunc("/{pid}", plant_controller.Delete).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.STATIC_FILES))))
	r.Use(loggingMw)
	r.Use(gardenMw)
	r.Use(logResultMw)
	return r
}

// loggingMw tell what route has been called.
func loggingMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Web: Received request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// gardenMw provides a garden for a given gid
func gardenMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if tools.HasValue(req, "gid") {
			log.Printf("Web: Request need a Garden")
			id, err := tools.GetInt(req, "gid")
			if err != nil {
				log.Printf("Web: No ID Provided, can't fetch %d", id)
				if tools.IsAPI(req) {
					tools.GenerateAPIError(w, "Invalid ID provided.")
				}

				return
			}
			handler := db.New()

			gard, err := garden.ByID(handler, id)

			if err != nil {
				log.Printf("Web: Failed to fetch Garden %d", id)
				if tools.IsAPI(req) {
					tools.GenerateAPIError(w, "Failed to fetch Garden")
				}

				return
			}

			context.Set(req, "garden", gard)
		}

		next.ServeHTTP(w, req)
	})
}

// loggingMw tell what route has been called.
func logResultMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, req)
		log.Printf(
			"%s %s (code=%d dt=%s written=%d)",
			req.Method,
			req.URL,
			m.Code,
			m.Duration,
			m.Written,
		)
	})
}

// ListenAndServe start listing http server
func ListenAndServe(router *mux.Router) {
	log.Printf("Web: Started server on 127.0.0.1:80 and listening ... ")
	http.ListenAndServe(":80", router)
}
