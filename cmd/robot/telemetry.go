package main

import (
	"context"
	"log"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

func runTelemetry(client pb.TelemetryServiceClient, robotID string) {
	stream, err := client.TelemetryStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		// if rand.Intn(100) < 5 {
		// 	log.Fatal("robot crashed during telemetry")
		// }

		err := stream.Send(&pb.TelemetryMessage{
			Version:         "v1",
			RobotId:         robotID,
			TimestampUnixMs: time.Now().UnixMilli(),
			Data:            []byte("telemetry"),
		})
		if err != nil {
			log.Println("telemetry failed:", err)
			time.Sleep(time.Second)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
