package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "jiuchen1986/hello-world/pkg/nettest"

	"google.golang.org/grpc"
)

// SimpleServer implements nettest grpc server
type SimpleServer struct {
	pb.UnimplementedNetTestServer
	Interval int
	Host     string
	Port     int
}

// NewServer returns an instance of SimpleServer
func NewServer(inv, port int, host string) *SimpleServer {
	return &SimpleServer{
		Interval: inv,
		Host:     host,
		Port:     port,
	}
}

// ListNumbers implements ListNumbers rpc of netest grpc server
func (ss *SimpleServer) ListNumbers(start *pb.Start, stream pb.NetTest_ListNumbersServer) error {
	log.Printf("get a request for numbers with start %d and length %d", start.GetNumber(), start.GetLength())
	sn, l := start.Number, start.Length

	for i := sn; i < sn+l; i++ {
		if err := stream.Send(&pb.Number{Number: i}); err != nil {
			log.Printf("failed to return number: %+v\n", err)
			return err
		}
		time.Sleep(time.Duration(ss.Interval) * time.Second)
	}
	return nil
}

// Run runs a nettest grpc server
func (ss *SimpleServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ss.Host, ss.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNetTestServer(grpcServer, ss)
	log.Printf("starting server listening at %s:%d", ss.Host, ss.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
