package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

func runTelemetry(client pb.TelemetryServiceClient, robotID string) {
	stream, err := client.TelemetryStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		// random crash
		if rand.Intn(100) < 5 {
			log.Fatal("robot crashed during telemetry")
		}

		msg := &pb.TelemetryMessage{
			Version:         "v1",
			RobotId:         robotID,
			TimestampUnixMs: time.Now().UnixMilli(),
			Data:            []byte("telemetry"),
		}

		if err := stream.Send(msg); err != nil {
			log.Println("telemetry send failed:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		time.Sleep(500 * time.Millisecond)
	}
}
