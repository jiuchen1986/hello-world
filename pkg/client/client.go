package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "jiuchen1986/hello-world/pkg/nettest"

	"google.golang.org/grpc"
)

// SimpleClient wraps a nettest grpc client
type SimpleClient struct {
	Client       pb.NetTestClient
	StartNumber  int32
	ReturnLength int32
	TestTimes    int
	Host         string
	Port         int
}

// NewClient returns an instance of SimpleClient
func NewClient(start, length int32, times, port int, host string) *SimpleClient {
	return &SimpleClient{
		StartNumber:  start,
		ReturnLength: length,
		TestTimes:    times,
		Host:         host,
		Port:         port,
	}
}

// netTest starts a net test, and exits on fatal error
func (sc *SimpleClient) netTest(pctx context.Context) {
	log.Printf("start net test with start number: %d, length: %d, times: %d\n",
		sc.StartNumber, sc.ReturnLength, sc.TestTimes)
	st := &pb.Start{
		Number: sc.StartNumber,
		Length: sc.ReturnLength,
	}
	cancels := []context.CancelFunc{}
	defer func() {
		for _, c := range cancels {
			c()
		}
	}()
	for i := 0; i < sc.TestTimes; i++ {
		ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
		cancels = append(cancels, cancel)
		stream, err := sc.Client.ListNumbers(ctx, st)
		if err != nil {
			log.Fatalf("failed to list numbers %+v at iteration %d\n", err, i)
		}

		for {
			number, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to list numbers %+v at iteration %d\n", err, i)
			}
			log.Printf("get number %d at iteration %d", number.Number, i)
		}
	}
}

// Run runs a nettest grpc client which continue doing net test
func (sc *SimpleClient) Run() {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", sc.Host, sc.Port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	sc.Client = pb.NewNetTestClient(conn)
	sc.netTest(context.Background())
}
