package main

import (
	"context"
	"github.com/r00takaspin/raft/grpc_api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

const (
	port = ":24816"
)

type Status byte

const (
	FOLLOWER Status = iota
	CANDIDATE
	LEADER
)

func (s *server) SetValue(ctx context.Context, r *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	log.Printf("SetValue: %v", r.Value)
	return &pb.SetValueResponse{Message: pb.ResponseStatus_OK}, nil
}

func (s *server) Init(ctx context.Context, r *pb.InitRequest) (*pb.InitResponse, error) {
	log.Printf("Init: %v", r.Nodes)
	return &pb.InitResponse{Message: pb.ResponseStatus_OK}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", port)
	}

	log.Printf("Listening grpc requests on port %v", port)
	s := grpc.NewServer()
	// Register reflection service on gRPC server.
	pb.RegisterRaftServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
