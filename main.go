package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/r00takaspin/raft/grpc_api/raft"
	"github.com/r00takaspin/raft/lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"time"
)

type server struct{}

type Status byte

const (
	MIN_HEARTBEAT = 150
	MAX_HEARTBEAT = 300
)

const (
	FOLLOWER Status = iota
	CANDIDATE
	LEADER
)

// Hostname and port of server
var Hostname string
var Port string

// Current value
var Value int32

// Address of leader node
var Leader pb.RaftServiceClient = nil

// Number of node votes
var Term = 0

//Current status of node
var NodeStatus Status

type Host struct {
	Hostname string
	Port     int
}

func (h *Host) toString() string {
	return fmt.Sprintf("%v:%v", h.Hostname, h.Port)
}

var host Host

// LogEntities
var LogNumber int32 = 0

type LogEntity struct {
	Id    int32
	Value int32
}

type LogEntities []LogEntity

func (logs *LogEntities) addLog(id int32, value int32) {
	*logs = append(*logs, LogEntity{Id: id, Value: value})
}

//Logs
var Logs LogEntities

//Node list for RPC calls
var Nodes []pb.RaftServiceClient

// gRPC
func (s *server) SetValue(ctx context.Context, r *pb.SetValueRequest) (*pb.StatusResponse, error) {
	if NodeStatus == LEADER {
		log.Printf("LEADER(%v) value from to %v", host.toString(), r.Value)

		LogNumber := LogNumber + 1
		Logs.addLog(LogNumber, r.Value)

		for i := 0; i < len(Nodes); i++ {
			Nodes[i].SetValue(ctx, &pb.SetValueRequest{LogId: int32(LogNumber), Value: r.Value})
		}
	} else {
		log.Printf("FOLLOVER(%v) changing value from to %v", Value, r.Value)

		Logs.addLog(r.LogId, r.Value)
	}
	return &pb.StatusResponse{Message: true}, nil
}

func (s *server) Heartbeat(ctx context.Context, r *pb.EmptyRequest) (*pb.StatusResponse, error) {
	//log.Printf("%v: Heartbeat", host.toString())
	return &pb.StatusResponse{Message: true}, nil
}

func (s *server) RequestVote(ctx context.Context, r *pb.EmptyRequest) (*pb.StatusResponse, error) {
	//log.Printf("%v: RequestVote", host.toString())
	return &pb.StatusResponse{Message: true}, nil
}

func (s *server) IsLeader(ctx context.Context, r *pb.EmptyRequest) (*pb.StatusResponse, error) {
	//log.Printf("%v: IsLeader: %v", host.toString(), NodeStatus == LEADER)

	if NodeStatus == LEADER {
		return &pb.StatusResponse{Message: true}, nil
	} else {
		return &pb.StatusResponse{Message: false}, nil
	}
}

func heartbeat(ctx context.Context, timeout int, nodes []pb.RaftServiceClient) {
	sleep(timeout)

	//log.Printf("Heartbeat from %v", host.toString())

	if Leader == nil {
		result := findLeader(ctx, nodes)
		if result == nil && NodeStatus != LEADER {
			makeElection(ctx, nodes)
		}
	} else {
		r, err := Leader.IsLeader(ctx, &pb.EmptyRequest{})
		if r.Message == false || err != nil {
			Leader = nil
			makeElection(ctx, nodes)
		}
	}
	heartbeat(ctx, timeout, nodes)
}

func findLeader(ctx context.Context, nodes []pb.RaftServiceClient) pb.RaftServiceClient {
	if NodeStatus == LEADER {
		return nil
	}

	//log.Printf("%v: finding leader", host.toString())

	for i := 0; i < len(nodes); i++ {
		currentNode := nodes[i]

		r, err := currentNode.IsLeader(ctx, &pb.EmptyRequest{})
		if err != nil {
			log.Printf("%v: IsLeader error: %v", host.toString(), err)
			continue
		}

		if r.Message == true {
			return currentNode
		}
	}
	return nil
}

func sleep(timeout int) {
	timeToSleep := time.Duration(timeout) * time.Millisecond
	time.Sleep(timeToSleep)
}

func makeElection(ctx context.Context, nodes []pb.RaftServiceClient) {
	NodeStatus = CANDIDATE

	log.Printf("%v becomes candidate", host.toString())

	Term = 1

	for i := 0; i < len(nodes); i++ {
		if NodeStatus == LEADER {
			return
		}

		consensus := int(math.Floor(float64(len(nodes)/2))) + 1

		if Term >= consensus {
			log.Printf("%v becomes Leader with %v votes of %v needed", host.toString(), Term, consensus)

			NodeStatus = LEADER
			Term = 0
			return
		}

		if requestVote(ctx, nodes[i]) == true {
			Term += 1
		}
	}
}

func requestVote(ctx context.Context, currentNode pb.RaftServiceClient) bool {
	r, error := currentNode.RequestVote(ctx, &pb.EmptyRequest{})
	if error != nil {
		log.Printf("%v: RequestVote error: %v", host.toString(), error)
		return false
	}

	if r.Message == true {
		log.Printf("%v got vote from %v", host.toString(), currentNode)
		return true
	}

	return false
}

func init() {
	NodeStatus = FOLLOWER
	Hostname, _ = os.Hostname()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	rand.Seed(time.Now().UTC().UnixNano())
	Port := flag.Int("p", 24816, "-p=10001")
	nodeArg := flag.String("nodes", "", "-nodes=node2:10002,node2:10003")
	flag.Parse()
	nodeAddresses := raft.ParseNodes(*nodeArg)

	host = Host{Hostname: Hostname, Port: *Port}

	address := fmt.Sprintf("%v:%v", Hostname, *Port)

	log.Printf("%v started as Follower", host.toString())

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen: %v", Port)
	}

	log.Printf("Listening grpc requests on Port %v", *Port)
	s := grpc.NewServer()
	// Register reflection service on gRPC server.

	ctx := context.Background()

	nodes := connectToNodes(nodeAddresses)
	Nodes = nodes

	timeout := rand.Intn(MAX_HEARTBEAT-MIN_HEARTBEAT) + MIN_HEARTBEAT
	go heartbeat(ctx, timeout, nodes)

	pb.RegisterRaftServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func connectToNodes(nodes []string) []pb.RaftServiceClient {
	var result []pb.RaftServiceClient

	for i := 0; i < len(nodes); i++ {
		currentNode := nodes[i]

		Term = Term + 1

		conn, err := grpc.Dial(currentNode, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		c := pb.NewRaftServiceClient(conn)

		result = append(result, c)
	}
	log.Printf("connected to %v nodes from %v", len(result), host.toString())

	return result
}
