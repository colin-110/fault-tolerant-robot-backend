package grpc

import (
	"github.com/colin-110/fault-tolerant-robot-backend/internal/workers"
	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommandServer struct {
	pb.UnimplementedCommandServiceServer
	commandPool *workers.CommandWorkerPool
}

func NewCommandServer(pool *workers.CommandWorkerPool) *CommandServer {
	return &CommandServer{commandPool: pool}
}

func (s *CommandServer) CommandStream(
	stream pb.CommandService_CommandStreamServer,
) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		ok := s.commandPool.Enqueue(workers.CommandJob{
			CommandID: req.CommandId,
		})
		if !ok {
			return status.Error(codes.ResourceExhausted, "command queue full")
		}

		if err := stream.Send(&pb.CommandResponse{
			Version:   req.Version,
			CommandId: req.CommandId,
			State:     pb.CommandState_SENT,
		}); err != nil {
			return err
		}
	}
}
