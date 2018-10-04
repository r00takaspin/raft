package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/r00takaspin/raft/grpc_api/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
)

type server struct{}

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
var Timeount int8

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
	port := flag.Int("p", 24816, "-p=10001")
	nodeList := flag.String("nodes", "", "-nodes=node2:10002,node2:10003")

	Nodes = strings.Split(*nodeList, ",")

	flag.Parse()

	Hostname, _ = os.Hostname()

	address := fmt.Sprintf("%v:%v", Hostname, *port)

	Timeount := rand.
		initState()

	log.Printf("%v started as Follower", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen: %v", port)
	}

	log.Printf("Listening grpc requests on port %v", *port)
	s := grpc.NewServer()
	// Register reflection service on gRPC server.
	pb.RegisterRaftServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initState() {
	NodeStatus = FOLLOWER
}
