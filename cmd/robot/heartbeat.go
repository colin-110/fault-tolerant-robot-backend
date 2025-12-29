package main

import (
	"context"
	"log"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

func runHeartbeat(client pb.HeartbeatServiceClient, robotID string) {
	stream, err := client.Beat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		err := stream.Send(&pb.Heartbeat{
			Version: "v1",
			RobotId: robotID,
		})
		if err != nil {
			log.Println("heartbeat failed:", err)
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
	}
}
