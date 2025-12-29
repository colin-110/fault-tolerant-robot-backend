package grpc

import (
	"github.com/colin-110/fault-tolerant-robot-backend/internal/workers"
	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

type TelemetryServer struct {
	pb.UnimplementedTelemetryServiceServer
	pool *workers.TelemetryWorkerPool
}

func NewTelemetryServer(pool *workers.TelemetryWorkerPool) *TelemetryServer {
	return &TelemetryServer{pool: pool}
}

func (s *TelemetryServer) TelemetryStream(
	stream pb.TelemetryService_TelemetryStreamServer,
) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		s.pool.Enqueue(workers.TelemetryJob{
			RobotID: msg.RobotId,
			Payload: msg.Data,
		})
	}
}
