package main

import (
	"context"
	"github.com/r00takaspin/raft/grpc_api/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
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

var Value int32
var Nodes []string
var NodeStatus Status
var Hostname string

func (s *server) SetValue(ctx context.Context, r *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	log.Printf("Changing value from %v to %v", Value, r.Value)
	Value = r.Value
	return &pb.SetValueResponse{Message: pb.ResponseStatus_OK}, nil
}

func (s *server) Init(ctx context.Context, r *pb.InitRequest) (*pb.InitResponse, error) {
	log.Printf("Init: %v", r.Nodes)
	return &pb.InitResponse{Message: pb.ResponseStatus_OK}, nil
}

func main() {
	Hostname, _ = os.Hostname()

	ctx := context.Background()
	initState(ctx)

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

func initState(ctx context.Context) {
	log.Printf("%v started as Follower", Hostname)

	NodeStatus = FOLLOWER
}
