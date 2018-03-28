package gardens

import (
	"errors"
	"log"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
)

// GardenDescriptor Allow quick access to Garden without whole structure loading.
type GardenDescriptor struct {
	id   int
	name string
}

// All fetches all garden from DB.
func All(dbh *db.Handler) []*garden.Garden {
	var res []*garden.Garden

	rows := dbh.Exec("SELECT * FROM garden")
	for rows.Next() {
		res = append(res, garden.Create(rows))
	}
	defer rows.Close()
	return res
}

// AllIds Fetches all ids and name of garden from DB
func AllIds(dbh *db.Handler) []GardenDescriptor {
	var res []GardenDescriptor

	rows := dbh.Exec("SELECT id,name FROM garden")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		res = append(res, GardenDescriptor{id, name})
	}
	defer rows.Close()
	return res
}

// ByID Fetch garden with provided id.
func ByID(dbh *db.Handler, id int) (*garden.Garden, error) {
	var res *garden.Garden

	rows := dbh.Exec("SELECT * FROM garden WHERE garden_id='?'", id)
	if rows.Next() {
		res = garden.Create(rows)
	} else {
		log.Printf("Gardens: No matching garden found: %d", id)
		return nil, errors.New("Gardens: No match found")
	}

	defer rows.Close()

	return res, nil
}
