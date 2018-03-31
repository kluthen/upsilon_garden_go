// Package garden_controller stores handles for garden requests.
package garden_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/lib/gardens"
	"upsilon_garden_go/web/templates"
	"upsilon_garden_go/web/tools"

	"github.com/gorilla/context"
)

// Index GET: /gardens
func Index(w http.ResponseWriter, req *http.Request) {
	// tools.IsAPI(req)
	handler := db.New()

	data := gardens.AllIds(handler)

	if tools.IsAPI(req) {
		tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(data)
	} else {
		templates.RenderTemplate(w, "garden\\index", data)
	}
}

// Show GET: /gardens/:id
func Show(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)

	if tools.IsAPI(req) {
		tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(gard)
	}
}

// Create POST: /gardens; expect name.
func Create(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	handler := db.New()

	gard := garden.New()
	gard.Name = name
	err := gard.Repsert(handler)

	if err != nil {
		log.Print("GardenCtrl: Failed to Create Garden.")
		if tools.IsAPI(req) {
			tools.GenerateAPIError(w, "Failed to create Garden")
		}

		return
	}

	if tools.IsAPI(req) {
		tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(gard)
	}
}

// Update PUT: /gardens/:id
func Update(w http.ResponseWriter, req *http.Request) {
	// may only update name.
	name := req.FormValue("name")
	handler := db.New()
	gard := context.Get(req, "garden").(*garden.Garden)
	gard.Name = name
	rerr := gard.Repsert(handler)

	if rerr != nil {
		log.Print("GardenCtrl: Failed to Create Garden.")
		if tools.IsAPI(req) {
			tools.GenerateAPIError(w, "Failed to create Garden")
		}

		return
	}

	if tools.IsAPI(req) {
		tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(gard)
	}

}

// Delete DELETE: /gardens/:id
func Delete(w http.ResponseWriter, req *http.Request) {
	handler := db.New()
	gard := context.Get(req, "garden").(*garden.Garden)
	gard.Drop(handler)

	if tools.IsAPI(req) {
		repm := tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(repm)
	}
}

// ShowHydro GET: /gardens/:id/hydro/:parcel
func ShowHydro(w http.ResponseWriter, req *http.Request) {

}

// AddHydro POST: /gardens/:id/hydro/:parcel
func AddHydro(w http.ResponseWriter, req *http.Request) {

}
