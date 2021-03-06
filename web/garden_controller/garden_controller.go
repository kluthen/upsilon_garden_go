// Package garden_controller stores handles for garden requests.
package garden_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
	defer handler.Close()

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
	} else {
		templates.RenderTemplate(w, "garden\\show", gard)
	}
}

// Create POST: /gardens; expect name.
func Create(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	handler := db.New()
	defer handler.Close()

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
	} else {
		http.Redirect(w, req, fmt.Sprintf("/gardens/%d", gard.ID), http.StatusSeeOther)
	}
}

// Update PUT: /gardens/:id
func Update(w http.ResponseWriter, req *http.Request) {
	// may only update name.
	name := req.FormValue("name")
	handler := db.New()
	defer handler.Close()
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
	} else {
		http.Redirect(w, req, fmt.Sprintf("/gardens/%d", gard.ID), http.StatusSeeOther)
	}

}

// Delete DELETE: /gardens/:id
func Delete(w http.ResponseWriter, req *http.Request) {
	handler := db.New()
	defer handler.Close()
	gard := context.Get(req, "garden").(*garden.Garden)
	gard.Drop(handler)

	if tools.IsAPI(req) {
		repm := tools.GenerateAPIOk(w)
		json.NewEncoder(w).Encode(repm)
	} else {
		http.Redirect(w, req, "/gardens", http.StatusSeeOther)
	}
}

// ShowHydro GET: /api/gardens/:id/hydro/:parcel
func ShowHydro(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	pid, _ := tools.GetInt(req, "parcel")

	parcel := gard.ParcelAt(pid)
	repm := tools.GenerateAPIOk(w)
	repm["hydro"] = parcel.GetHumanCurrentHydro()

	json.NewEncoder(w).Encode(repm)
}

// AddHydro POST: /api/gardens/:id/hydro/:parcel
func AddHydro(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	pid, _ := tools.GetInt(req, "parcel")

	parcel := gard.ParcelAt(pid)
	dur, _ := time.ParseDuration("8h")
	parcel.AddHydroEvent(time.Now().Add(dur), 0.15)

	handler := db.New()
	defer handler.Close()
	gard.Repsert(handler)

	repm := tools.GenerateAPIOk(w)
	repm["hydro"] = parcel.GetHumanCurrentHydro()

	json.NewEncoder(w).Encode(repm)
}
