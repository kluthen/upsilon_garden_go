package main

import (
	"log"
	"time"
	"upsilon_garden_go/lib/db"
	"upsilon_garden_go/lib/garden"
	"upsilon_garden_go/lib/gardens"
	"upsilon_garden_go/web"
)

func testDB(handler *db.Handler) {

	gard := garden.New()
	gard.Name = "Kluthen's"
	gard.Repsert(handler)
	res := gardens.AllIds(handler)

	for _, row := range res {
		log.Printf("Test: Fetched: %s", row.String())
		gard, err := garden.ByID(handler, row.ID)
		if err != nil {
			log.Fatalf("Test: Failed to fetch garden : %d", row.ID)
		}
		log.Printf("Test: Fetched Garden: %s", gard.String())

		gard.LastUpdate = time.Now()
		gard.Repsert(handler)

		gard, err = garden.ByID(handler, row.ID)
		if err != nil {
			log.Fatalf("Test: Failed to fetch garden : %d", row.ID)
		}
		log.Printf("Test: Fetched Garden: %s", gard.String())

		s_garden := garden.New()
		s_garden.Name = "Test's"
		s_garden.Repsert(handler)

		var ids []int
		ids = append(ids, gard.ID)
		ids = append(ids, s_garden.ID)

		resGardens, err := garden.ByIDs(handler, ids)
		if err != nil {
			log.Fatalf("Test: Failed to find ids : %s", err)
		} else {
			if len(resGardens) != 2 {
				log.Fatalf("Test: Hasn't found appropriate number of rows: %d (expected 2)", len(resGardens))
			} else {
				for i := range resGardens {
					log.Printf("Test: Found in ids: %s", resGardens[i].String())
				}
			}
		}

		gard.Drop(handler)
		s_garden.Drop(handler)

		gard, err = garden.ByID(handler, row.ID)
		if err != nil {
			log.Printf("Test: Successfully lost : %d", row.ID)
		} else {
			log.Fatalf("Test: Found Garden that should have been dropped %d", row.ID)
		}
	}
}

func main() {
	handler := db.New()
	// testDB(handler)
	r := web.RouterSetup()
	web.ListenAndServe(r)

	defer handler.Close()

}
