package grpc

import (
	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
	pb "github.com/colin-110/fault-tolerant-robot-backend/proto"
)

type HeartbeatServer struct {
	pb.UnimplementedHeartbeatServiceServer
	store *storage.LivenessStore
}

func NewHeartbeatServer(store *storage.LivenessStore) *HeartbeatServer {
	return &HeartbeatServer{store: store}
}

func (s *HeartbeatServer) Beat(
	stream pb.HeartbeatService_BeatServer,
) error {
	for {
		hb, err := stream.Recv()
		if err != nil {
			return err
		}
		_ = s.store.Heartbeat(hb.RobotId)
	}
}
