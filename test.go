package main

import (
	"upsilon_garden_go/lib/db"
)

func main() {
	handler := db.New()
	defer handler.Close()

}
