package garden

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"upsilon_garden_go/lib/db"

	"github.com/lib/pq"
)

// Garden related data.
type Garden struct {
	ID         int
	Name       string
	LastUpdate time.Time
	NextUpdate time.Time
	Parcels    []Parcel
	Plants     []Plant
}

// New empty Garden
func New() *Garden {
	garden := new(Garden)
	garden.ID = 0
	garden.LastUpdate = time.Now()

	// force 10 parcels in a garden.

	for i := 0; i < 10; i++ {
		np := newParcel()
		np.Position = i
		garden.Parcels = append(garden.Parcels, np)
	}

	return garden
}

// String pretty displayer
func (garden *Garden) String() string {
	return fmt.Sprintf("Garden { "+
		"ID: %d"+
		", name: %s"+
		", lastUpdate: %s"+
		", nextUpdate: %s"+
		" }", garden.ID, garden.Name, garden.LastUpdate, garden.NextUpdate)
}

// Create a Garden from database
func Create(rows *sql.Rows) *Garden {
	garden := new(Garden)
	var parcels []byte
	var plants []byte
	rows.Scan(&garden.ID, &garden.Name, &garden.LastUpdate, &garden.NextUpdate, &parcels, &plants)
	json.Unmarshal(parcels, &garden.Parcels)
	json.Unmarshal(plants, &garden.Plants)
	return garden
}

// Repsert Garden in Database
func (garden *Garden) Repsert(dbh *db.Handler) error {
	if garden.ID <= 0 {
		// Need to insert it at first.
		rows := dbh.Exec("INSERT INTO garden(name) VALUES (NULL) RETURNING garden_id")
		for rows.Next() {
			rows.Scan(&garden.ID)
		}

		if garden.ID <= 0 {
			log.Fatal("Garden: Failed to Insert garden in database !")
			return errors.New("Unable to create garden row")
		}
		rows.Close()
	}

	parcels, err := json.Marshal(garden.Parcels)
	if err != nil {
		log.Fatal("Garden: Failed to jsonify Parcels")
		return errors.New("Unable to jsonify Parcels")
	}

	plants, err := json.Marshal(garden.Plants)
	if err != nil {
		log.Fatal("Garden: Failed to jsonify Plants")
		return errors.New("Unable to jsonify Plants")
	}

	rows := dbh.Query(`UPDATE garden SET 
			name=$1 ,
			last_update=$2 ,
			next_update=$3 ,
			parcels=$4 ,
			plants=$5
			WHERE garden_ID=$6 
		`, garden.Name,
		garden.LastUpdate,
		garden.NextUpdate,
		parcels,
		plants,
		garden.ID)

	rows.Close()

	return nil
}

// Drop Garden from database
func (garden *Garden) Drop(dbh *db.Handler) {
	if garden.ID <= 0 {
		log.Print("Garden: Can't drop a non existing Garden.")
		return
	}

	rows := dbh.Query("DELETE FROM garden WHERE garden_ID=$1", garden.ID)
	rows.Close()

	garden.ID = 0
}

// ByID Fetch garden with provided id.
func ByID(dbh *db.Handler, id int) (*Garden, error) {
	var res *Garden

	rows := dbh.Query("SELECT garden_id, name, last_update, next_update, parcels, plants FROM garden WHERE garden_id=$1", id)
	if rows.Next() {
		res = Create(rows)
	} else {
		log.Printf("Gardens: No matching garden found: %d", id)
		return nil, errors.New("Gardens: No match found")
	}

	defer rows.Close()

	return res, nil
}

// ByIDs Fetch garden with provided id.
func ByIDs(dbh *db.Handler, ids []int) ([]*Garden, error) {
	var res *Garden
	var results []*Garden

	if len(ids) == 0 {
		log.Printf("Gardens: Hasn't provided any id to fetch")
		return nil, errors.New("Gardens: No target provided")
	}

	rows := dbh.Query("SELECT garden_id, name, last_update, next_update, parcels, plants FROM garden WHERE garden_id = ANY($1)", pq.Array(ids))
	for rows.Next() {
		res = Create(rows)
		results = append(results, res)
	}
	if len(results) == 0 {
		log.Printf("Gardens: No matching garden found: %v", ids)
		return nil, errors.New("Gardens: No match found")
	}

	defer rows.Close()

	return results, nil
}
