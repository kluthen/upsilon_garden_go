package garden

import (
	"time"
)

// HydroEvent stores information related to Hydro events duration and power.
type HydroEvent struct {
	BeginDate time.Time
	EndDate   time.Time
	Power     float32
}
