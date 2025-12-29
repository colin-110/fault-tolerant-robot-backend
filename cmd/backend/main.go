package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
	"github.com/colin-110/fault-tolerant-robot-backend/internal/workers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := storage.OpenDB("robot.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	commandStore := storage.NewCommandStore(db)

	commandPool := workers.NewCommandWorkerPool(commandStore)
	telemetryPool := workers.NewTelemetryWorkerPool()

	commandPool.Start(ctx)
	telemetryPool.Start(ctx)

	log.Println("backend started with worker pools")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutdown signal received")
	cancel()
	log.Println("backend shutting down")
}
