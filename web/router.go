package web

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"upsilon_garden_go/config"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/web/garden_controller"
	"upsilon_garden_go/web/plant_controller"
	"upsilon_garden_go/web/templates"
	"upsilon_garden_go/web/tools"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// RouterSetup Prepare routing.
func RouterSetup() *mux.Router {
	rand.Seed(time.Now().Unix())
	r := mux.NewRouter()

	// Hydro functions
	r.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	r.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")
	// CRUD /gardens
	r.HandleFunc("/gardens/{gid}", garden_controller.Show).Methods("GET")
	r.HandleFunc("/gardens/{gid}", garden_controller.Update).Methods("PUT")
	r.HandleFunc("/gardens/{gid}", garden_controller.Delete).Methods("DELETE")
	r.HandleFunc("/gardens", garden_controller.Index).Methods("GET")
	r.HandleFunc("/gardens", garden_controller.Create).Methods("POST")

	// CRUD /gardens/:id/plant
	plantRouter := r.PathPrefix("/gardens/{gid}").Subrouter()
	plantRouter.HandleFunc("/plants/{pid}", plant_controller.Show).Methods("GET")
	plantRouter.HandleFunc("/plants/{pid}", plant_controller.Update).Methods("PUT")
	plantRouter.HandleFunc("/plants/{pid}", plant_controller.Delete).Methods("DELETE")
	plantRouter.HandleFunc("/plants/", plant_controller.Index).Methods("GET")
	plantRouter.HandleFunc("/plants/", plant_controller.Create).Methods("POST")

	// JSON Access ...

	jsonAPI := r.PathPrefix("/api").Subrouter()

	// Hydro functions
	jsonAPI.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.ShowHydro).Methods("GET")
	jsonAPI.HandleFunc("/gardens/{gid}/hydro/{parcel}", garden_controller.AddHydro).Methods("POST")

	// CRUD /gardens
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Show).Methods("GET")
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Update).Methods("PUT")
	jsonAPI.HandleFunc("/gardens/{gid}", garden_controller.Delete).Methods("DELETE")
	jsonAPI.HandleFunc("/gardens", garden_controller.Index).Methods("GET")
	jsonAPI.HandleFunc("/gardens", garden_controller.Create).Methods("POST")

	// CRUD /gardens/:id/plant
	APIPlantRouter := jsonAPI.PathPrefix("/gardens/{gid}").Subrouter()
	APIPlantRouter.HandleFunc("/plants/{pid}", plant_controller.Show).Methods("GET")
	APIPlantRouter.HandleFunc("/plants/{pid}", plant_controller.Update).Methods("PUT")
	APIPlantRouter.HandleFunc("/plants/{pid}", plant_controller.Delete).Methods("DELETE")
	APIPlantRouter.HandleFunc("/plants/", plant_controller.Index).Methods("GET")
	APIPlantRouter.HandleFunc("/plants/", plant_controller.Create).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.STATIC_FILES))))

	r.Use(logResultMw)
	r.Use(loggingMw)
	r.Use(gardenMw)
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
			defer handler.Close()

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
			"Web: %s %s (code=%d dt=%s written=%d)",
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
	templates.LoadTemplates()
	log.Printf("Web: Started server on 127.0.0.1:80 and listening ... ")
	http.ListenAndServe(":80", router)
}
