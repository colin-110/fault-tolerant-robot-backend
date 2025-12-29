package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	robotID := os.Getenv("ROBOT_ID")
	if robotID == "" {
		log.Fatal("ROBOT_ID not set")
	}

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	commandClient := pb.NewCommandServiceClient(conn)
	telemetryClient := pb.NewTelemetryServiceClient(conn)
	heartbeatClient := pb.NewHeartbeatServiceClient(conn)

	go runTelemetry(telemetryClient, robotID)
	go runHeartbeat(heartbeatClient, robotID)
	runCommands(commandClient, robotID)
}
