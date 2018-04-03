package garden

import (
	"fmt"
	"log"
	"time"
	"upsilon_garden_go/config"
)

// Plant stores plant related information, like, name, type, events.
type Plant struct {
	ID          int
	Name        string
	Level       int
	PlantType   string
	TargetHydro HydroRange
	NextUpdate  time.Time
	SpPerLevel  int
	SpMax       int
	SpCurrent   int
}

// NewPlant Generate a new plant.
func NewPlant() Plant {
	var plant Plant
	plant.Name = "Some Plant"
	plant.Level = 1
	plant.PlantType = "Some Type"
	plant.TargetHydro.MaxNotDead = 0.85
	plant.TargetHydro.MinNotDead = 0.25
	plant.TargetHydro.MaxOk = 0.75
	plant.TargetHydro.MinOk = 0.35
	plant.TargetHydro.MaxSuper = 0.65
	plant.TargetHydro.MinSuper = 0.45
	dur, _ := time.ParseDuration("3h")
	plant.NextUpdate = time.Now().UTC().Add(dur)
	plant.SpPerLevel = 5
	plant.SpMax = 50
	plant.SpCurrent = 25

	return plant
}

func (plant *Plant) String() string {
	if plant == nil {
		return "Nil Plant"
	}
	return fmt.Sprintf("Plant {"+
		"ID: %d"+
		", Name: %s"+
		", Level: %d"+
		", PlantType: %s"+
		", SP: %d"+
		", Max SP: %d",
		plant.ID, plant.Name, plant.Level, plant.PlantType, plant.SpCurrent, plant.SpMax)
}

// checkAndUpdate check alteration to the plant (sp, level up,etc) update stats accordingly and tell if plant has been altered.
// returns ()
func (plant *Plant) checkAndUpdate(lastVisit time.Time, toDate time.Time, currentHydro float32) (altered bool, destroyed bool) {
	altered = false
	destroyed = false
	modifier := plant.TargetHydro.WhereInRange(currentHydro)
	log.Printf("Plant: About to update: Plant NextUpdate: %v, Last Visit: %v, Target Date %v", plant.NextUpdate, lastVisit, toDate)
	// there was at least a Level up inbetween ...
	if toDate.Sub(plant.NextUpdate).Seconds() > 0 {
		timeSpent := int(toDate.Sub(plant.NextUpdate).Seconds()) / config.PLANT_TIME_TIC
		log.Printf("Plant: %d Tic to apply to plant %d", timeSpent, plant.ID)
		if timeSpent > 0 {
			altered = true
			plant.SpCurrent += modifier * timeSpent
			if plant.SpCurrent <= 0 {
				destroyed = true
				return
			}
			if plant.SpCurrent > plant.SpMax {
				plant.SpCurrent = plant.SpMax
			}

			lastUpdate := plant.NextUpdate
			plant.LevelUp()

			// There might have been more than one level up inbetween ;)
			return plant.checkAndUpdate(lastUpdate, toDate, currentHydro)
		}
		return
	}

	timeSpent := int(toDate.Sub(lastVisit).Seconds()) / config.PLANT_TIME_TIC
	log.Printf("Plant: %d Tic to apply to plant %d", timeSpent, plant.ID)

	if timeSpent > 0 {
		altered = true
		plant.SpCurrent += modifier * timeSpent
		if plant.SpCurrent <= 0 {
			destroyed = true
			return
		}
		if plant.SpCurrent > plant.SpMax {
			plant.SpCurrent = plant.SpMax
		}
	}

	return
}

// LevelUp levels a plant up (ahah)
func (plant *Plant) LevelUp() {
	plant.Level++
	plant.SpCurrent += plant.SpPerLevel
	plant.SpMax += plant.SpPerLevel
	dur, _ := time.ParseDuration("3h")
	for i := 0; i < plant.Level; i++ {
		plant.NextUpdate = plant.NextUpdate.Add(dur)
	}
}
