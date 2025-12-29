package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

func runCommands(client pb.CommandServiceClient, robotID string) {
	stream, err := client.CommandStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		// simulate reconnect / crash
		if rand.Intn(100) < 5 {
			log.Fatal("robot crashed during command handling")
		}

		resp, err := stream.Recv()
		if err != nil {
			log.Println("command stream error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("robot %s received command %s state=%s",
			robotID,
			resp.CommandId,
			resp.State,
		)
	}
}
