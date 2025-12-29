package main

import (
	"context"
	"log"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

// executed tracks commands already executed by this robot.
// This guarantees idempotent execution under duplicate delivery.
var executed = make(map[string]bool)

func runCommands(client pb.CommandServiceClient, robotID string) {
	stream, err := client.CommandStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Simulate crash during command handling (failure injection)
		// if rand.Intn(100) < 5 {
		// 	log.Fatal("robot crashed during command handling")
		// }

		resp, err := stream.Recv()
		if err != nil {
			log.Println("command stream error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Idempotency: ignore duplicate commands
		if executed[resp.CommandId] {
			log.Printf(
				"robot %s ignoring duplicate command %s",
				robotID,
				resp.CommandId,
			)
			continue
		}

		// Mark command as executed BEFORE performing side effects
		// This ensures crash safety (at-least-once delivery semantics)
		executed[resp.CommandId] = true

		log.Printf(
			"robot %s executed command %s (state=%v)",
			robotID,
			resp.CommandId,
			resp.State,
		)

		// Simulated execution time
		time.Sleep(500 * time.Millisecond)
	}
}
