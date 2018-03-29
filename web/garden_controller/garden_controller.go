// Package garden_controller stores handles for garden requests.
package garden_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/lib/gardens"
	"upsilon_garden_go/web/tools"
)

// Index GET: /gardens
func Index(w http.ResponseWriter, req *http.Request) {
	// tools.IsAPI(req)
	handler := db.New()

	if tools.IsAPI(req) {
		json.NewEncoder(w).encode(gardens.AllIds(handler))
	}
}

// Show GET: /gardens/:id
func Show(w http.ResponseWriter, req *http.Request) {
	id := tools.GetInt(req, "id")
	handler := db.New()

	garden, err = garden.ByID(handler, id)

	if err != nil {
		log.Printf("GardenCtrl: Failed to fetch Garden %d", id)
		if tools.IsAPI(req) {
			tools.GenerateAPIError(&w, "Failed to fetch Garden")
		}

		return
	}

	if tools.IsAPI(req) {
		json.NewEncoder(w).encode(garden)
	}
}

// Create POST: /gardens; expect name.
func Create(w http.ResponseWriter, req *http.Request) {
	name := r.FormValue("name")
	handler := db.New()

	garden := garden.New()
	garden.Name = name
	err = garden.Repsert(handler)

	if err != nil {
		log.Print("GardenCtrl: Failed to Create Garden.")
		if tools.IsAPI(req) {
			tools.GenerateAPIError(&w, "Failed to create Garden")
		}

		return
	}

	if tools.IsAPI(req) {
		json.NewEncoder(w).encode(garden)
	}
}

// Update PUT: /gardens/:id
func Update(rep http.ResponseWriter, req *http.Request) {

}

// Delete DELETE: /gardens/:id
func Delete(rep http.ResponseWriter, req *http.Request) {

}

// ShowHydro GET: /gardens/:id/hydro/:parcel
func ShowHydro(rep http.ResponseWriter, req *http.Request) {

}

// AddHydro POST: /gardens/:id/hydro/:parcel
func AddHydro(rep http.ResponseWriter, req *http.Request) {

}
