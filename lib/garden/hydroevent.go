package garden

import (
	"fmt"
	"time"
)

// HydroEvent stores information related to Hydro events duration and power.
type HydroEvent struct {
	BeginDate time.Time
	EndDate   time.Time
	Power     float32
}

func (hevt *HydroEvent) String() string {
	return fmt.Sprintf("HydroEvent {"+
		"Begin: %v, End: %v, Power: %f}",
		hevt.BeginDate, hevt.EndDate, hevt.Power)
}
