package garden

import (
	"database/sql"
	"time"
)

// Garden related data.
type Garden struct {
	id         int
	name       string
	lastUpdate time.Time
	nextUpdate time.Time
	parcels    []Parcel
	plants     []Plant
}

// New empty Garden
func New() *Garden {
	return new(Garden)
}

// Create a Garden from database
func Create(rows *sql.Rows) *Garden {
	garden := new(Garden)

	return garden
}

// Repsert Garden in Database
func (garden *Garden) Repsert() {

}

// Drop Garden from database
func (garden *Garden) Drop() {

}
