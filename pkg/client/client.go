package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "jiuchen1986/hello-world/pkg/nettest"
	"jiuchen1986/hello-world/pkg/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SimpleClient wraps a nettest grpc client
type SimpleClient struct {
	Client       pb.NetTestClient
	StartNumber  int32
	ReturnLength int32
	TestTimes    int
	Host         string
	Port         int
	Timeout      int
	CaCert       string
	Cert         string
	Key          string
}

// NewClient returns an instance of SimpleClient
func NewClient(start, length int32, times, port, 
                   timeout int, host string, cacert, cert, key string) *SimpleClient {
	return &SimpleClient{
		StartNumber:  start,
		ReturnLength: length,
		TestTimes:    times,
		Host:         host,
		Port:         port,
		Timeout:      timeout,
		CaCert:       cacert,
		Cert:         cert,
		Key:          key,
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
		ctx, cancel := context.WithTimeout(pctx, time.Duration(sc.Timeout)*time.Second)
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
	if sc.CaCert != "" {
	       tlsConfig, err := util.GetClientTLSConfig(sc.CaCert, sc.Cert, sc.Key, sc.Host)
	       if err != nil {
	               log.Fatalf("fail to create tls: %v", err)
               }
	       dialOption := grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	       opts = []grpc.DialOption{dialOption}
	}
	log.Printf("connecting %s:%d", sc.Host, sc.Port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", sc.Host, sc.Port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	sc.Client = pb.NewNetTestClient(conn)
	sc.netTest(context.Background())
}
