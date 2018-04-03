// Package plant_controller stores handles for plant requests.
package plant_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/web/tools"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Index GET: /gardens/:gid/plants
func Index(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)

	if tools.IsAPI(req) {
		data, err := json.Marshal(gard.Plants)
		if err != nil {
			log.Printf("PlantCtrl: Failed to Serialize plants of garden %d", gard.ID)
			tools.GenerateAPIError(w, "Failed to prepare reply.")
			return
		}

		repm := tools.GenerateAPIOk(w)
		repm["plants"] = string(data)
		json.NewEncoder(w).Encode(repm)
	}
}

// Show GET: /gardens/:gid/plants/:pid
func Show(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	plant := context.Get(req, "plant").(*garden.Plant)

	if tools.IsAPI(req) {
		data, err := json.Marshal(plant)
		if err != nil {
			log.Printf("PlantCtrl: Failed to Serialize plants of garden %d", gard.ID)
			tools.GenerateAPIError(w, "Failed to prepare reply.")
			return
		}

		repm := tools.GenerateAPIOk(w)
		repm["plant"] = string(data)
		json.NewEncoder(w).Encode(repm)
	}
}

// Create POST: /gardens/:gid/plants
// expect name, parcel, plant_type
func Create(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	plant := garden.NewPlant()
	plant.Name = req.FormValue("name")
	plant.PlantType = req.FormValue("plant_type")
	pid, err := strconv.Atoi(req.FormValue("parcel"))
	if err != nil {
		if tools.IsAPI(req) {
			log.Printf("PlantCtrl: Failed to read target parcel: %s", err)
			tools.GenerateAPIError(w, "Failed to prepare reply.")
		}
		return
	}
	err = gard.AddPlant(pid, plant)

	if err != nil {
		if tools.IsAPI(req) {
			log.Printf("PlantCtrl: Failed to create plants of garden %d: %s", gard.ID, err)
			tools.GenerateAPIError(w, "Failed to prepare reply.")
			return
		}
	} else {
		log.Printf("PlantCtrl: Successfully added a Plant at parcel: %d", pid)
	}

	if tools.IsAPI(req) {
		tools.GenerateAPIOkAndSend(w)
	}
}

// Update PUT: /gardens/:gid/plants/:pid
func Update(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	plant := context.Get(req, "plant").(*garden.Plant)
	vars := mux.Vars(req)
	plant.Name = vars["name"]

	handler := db.New()
	defer handler.Close()
	gard.Repsert(handler)

	if tools.IsAPI(req) {
		tools.GenerateAPIOkAndSend(w)
	}
}

// Delete DELETE: /gardens/:gid/plants/:pid
func Delete(w http.ResponseWriter, req *http.Request) {
	gard := context.Get(req, "garden").(*garden.Garden)
	plant := context.Get(req, "plant").(*garden.Plant)

	handler := db.New()
	defer handler.Close()
	gard.DropPlant(plant.ID)
	gard.Repsert(handler)

	if tools.IsAPI(req) {
		tools.GenerateAPIOkAndSend(w)
	}
}
