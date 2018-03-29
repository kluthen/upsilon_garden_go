package gardens

import (
	"fmt"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
)

// GardenDescriptor Allow quick access to Garden without whole structure loading.
type GardenDescriptor struct {
	ID   int
	Name string
}

// String pretty displayer
func (gd *GardenDescriptor) String() string {
	return fmt.Sprintf("GardenDescriptor { "+
		"ID: %d"+
		", name: %s }", gd.ID, gd.Name)
}

// All fetches all garden from DB.
func All(dbh *db.Handler) []*garden.Garden {
	var res []*garden.Garden

	rows := dbh.Exec("SELECT garden_id, name, last_update, next_update, parcels, plants FROM garden")
	for rows.Next() {
		res = append(res, garden.Create(rows))
	}
	defer rows.Close()
	return res
}

// AllIds Fetches all ids and name of garden from DB
func AllIds(dbh *db.Handler) []GardenDescriptor {
	var res []GardenDescriptor

	rows := dbh.Exec("SELECT garden_id,name FROM garden")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		res = append(res, GardenDescriptor{id, name})
	}
	defer rows.Close()
	return res
}
