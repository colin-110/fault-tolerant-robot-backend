package main

import (
	"log"

	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
)

func main() {
	db, err := storage.OpenDB("robot.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	log.Println("backend started, database initialized")
}
