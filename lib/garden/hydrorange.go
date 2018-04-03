package garden

import "upsilon_garden_go/config"

// HydroRange tells in what hydro level the plant may grow or not.
type HydroRange struct {
	MinNotDead float32
	MaxNotDead float32
	MinOk      float32
	MaxOk      float32
	MinSuper   float32
	MaxSuper   float32
}

const (
	HR_DD  int = config.PLANT_DD_TIC
	HR_ND  int = config.PLANT_ND_TIC
	HR_OK  int = config.PLANT_OK_TIC
	HR_SUP int = config.PLANT_SUP_TIC
)

// WhereInRange checks an hydro level against provided hydro range and tell where it lies
func (hr *HydroRange) WhereInRange(hydroLevel float32) int {
	if hydroLevel < hr.MinNotDead || hydroLevel > hr.MaxNotDead {
		return HR_DD
	}
	if hydroLevel < hr.MinOk || hydroLevel > hr.MaxOk {
		return HR_ND
	}
	if hydroLevel < hr.MinSuper || hydroLevel > hr.MaxSuper {
		return HR_OK
	}
	return HR_SUP

}
