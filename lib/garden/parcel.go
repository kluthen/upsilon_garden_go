package garden

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Parcel tell parcel related information like hydrometry, if it has a plant and so on.
type Parcel struct {
	ID                 int
	Position           int
	CurrentHydroLevel  float32
	BaseHydroLevel     float32
	PlantID            int
	RunningHydroEvents []HydroEvent
	NextHydroEnd       time.Time
}

func (parcel *Parcel) String() string {
	var events string

	for idx := range parcel.RunningHydroEvents {
		events += parcel.RunningHydroEvents[idx].String() + " "
	}

	return fmt.Sprintf("Parcel {"+
		"ID: %d"+
		", Position: %d"+
		", CurrentHydroLevel: %f"+
		", BaseHydroLevel: %f"+
		", PlantID: %d"+
		", Events: %s"+
		", NextHydroEnd: %v }",
		parcel.ID, parcel.Position,
		parcel.CurrentHydroLevel, parcel.BaseHydroLevel,
		parcel.PlantID, events, parcel.NextHydroEnd)
}

var baseHydroAvailable = [4]float32{0.15, 0.2, 0.25, 0.30}

// newParcel create a new random parcel.
func newParcel() Parcel {
	var p Parcel
	p.ID = 0
	p.Position = 0
	p.BaseHydroLevel = baseHydroAvailable[rand.Intn(len(baseHydroAvailable))]
	p.CurrentHydroLevel = p.BaseHydroLevel
	p.PlantID = -1
	p.NextHydroEnd = time.Time{}
	return p
}

// GetHumanCurrentHydro tell in a human (french) way how much  water there is on this plot.
func (parcel *Parcel) GetHumanCurrentHydro() string {
	if parcel.CurrentHydroLevel < 0.2 {
		return "dur"
	}
	if parcel.CurrentHydroLevel < 0.4 {
		return "sec"
	}
	if parcel.CurrentHydroLevel < 0.6 {
		return "normal"
	}
	if parcel.CurrentHydroLevel < 0.8 {
		return "humide"
	}
	if parcel.CurrentHydroLevel < 0.95 {
		return "boueux"
	}
	if parcel.CurrentHydroLevel < 1 {
		return "submergé"
	}
	return "normal"
}

// HasNextHydroEndDate tell whether an end date has been set or not.
func (parcel *Parcel) HasNextHydroEndDate() bool {
	return parcel.NextHydroEnd != time.Time{}
}

// GetNextHydroEndDate tell when next hydro event will end
func (parcel *Parcel) getNextHydroEndDate() (time.Time, bool) {
	t := time.Time{}
	var found bool
	found = false

	for idx := range parcel.RunningHydroEvents {
		event := &parcel.RunningHydroEvents[idx]
		if !found {
			t = event.EndDate
			found = true
		} else {
			if event.EndDate.Sub(t).Seconds() < 0 {
				t = event.EndDate
			}
		}

	}

	return t, found
}

// refreshParcel will check plant and parcel for update.
func (parcel *Parcel) refreshParcel(now time.Time, lastVisit time.Time, plant *Plant) (altered bool, plantDestroyed bool) {
	altered = false
	plantDestroyed = false

	if plant == nil {
		altered = parcel.checkAndRecomputeHydro()
		return
	}

	if parcel.HasNextHydroEndDate() {
		if now.Sub(parcel.NextHydroEnd).Seconds() < 0 {
			altered, plantDestroyed = plant.checkAndUpdate(lastVisit, now, parcel.CurrentHydroLevel)
			if altered {
				log.Printf("Parcel: Plant Altered: %s", plant.String())
			}
			return
		}
	} else {
		altered, plantDestroyed = plant.checkAndUpdate(lastVisit, now, parcel.CurrentHydroLevel)
		if altered {
			log.Printf("Parcel: Plant Altered: %s", plant.String())
		}
		return
	}

	altered = true
	lastNow := parcel.NextHydroEnd
	palter, plantDestroyed := plant.checkAndUpdate(lastVisit, parcel.NextHydroEnd, parcel.CurrentHydroLevel)

	if palter {
		log.Printf("Parcel: Plant Altered: %s", plant.String())
	}
	if plantDestroyed {
		parcel.checkAndRecomputeHydro()
		return
	}

	parcel.refreshParcel(now, lastNow, plant)
	return
}

// checkAndRecomputeHydro validate hydro events and remove them if appropriate. Tell if an alteration has been done.
func (parcel *Parcel) checkAndRecomputeHydro() bool {
	var addedPower float32
	var newHydroEvents []HydroEvent
	now := time.Now()
	altered := false

	addedPower += parcel.BaseHydroLevel
	for idx := range parcel.RunningHydroEvents {
		event := parcel.RunningHydroEvents[idx]
		if event.EndDate.Sub(now).Seconds() > 0 {
			addedPower += event.Power
			newHydroEvents = append(newHydroEvents, event)
		} else {
			altered = true
		}
	}

	if addedPower < 0.99 {
		parcel.CurrentHydroLevel = addedPower
	} else {
		parcel.CurrentHydroLevel = 0.99
	}

	parcel.RunningHydroEvents = newHydroEvents

	parcel.NextHydroEnd, _ = parcel.getNextHydroEndDate()

	return altered
}

// AddHydroEvent create and add hydro event to a parcel.
func (parcel *Parcel) AddHydroEvent(endDate time.Time, power float32) {
	var event HydroEvent
	event.BeginDate = time.Now()
	event.EndDate = endDate
	event.Power = power
	parcel.RunningHydroEvents = append(parcel.RunningHydroEvents, event)

	parcel.checkAndRecomputeHydro()
}
