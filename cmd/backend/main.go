package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "github.com/colin-110/fault-tolerant-robot-backend/internal/grpc"
	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
	"github.com/colin-110/fault-tolerant-robot-backend/internal/workers"
	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
	"google.golang.org/grpc"
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer()

	pb.RegisterCommandServiceServer(
		grpcSrv,
		grpcserver.NewCommandServer(commandPool),
	)
	pb.RegisterTelemetryServiceServer(
		grpcSrv,
		grpcserver.NewTelemetryServer(telemetryPool),
	)

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutdown initiated")
	cancel()
	grpcSrv.GracefulStop()
	log.Println("backend stopped")
}
