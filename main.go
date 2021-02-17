package main

import (
	"flag"
	"log"

	"jiuchen1986/hello-world/pkg/client"
	"jiuchen1986/hello-world/pkg/server"
)

var (
	startNumber  = flag.Int("start", 0, "start number that server returns since")
	returnLength = flag.Int("length", 5, "length of number sequence that server returns")
	testTimes    = flag.Int("times", 10, "number of times to test")
	interval     = flag.Int("interval", 1, "interval between returning numbers in seconds")
	host         = flag.String("host", "localhost", "host which server listening at, or client connecting to")
	port         = flag.Int("port", 5473, "port which server listening at, or client connecting to")
	mode         = flag.String("mode", "server", "support \"server\" and \"client\"")
)

func main() {

	flag.Parse()
	if *mode == "client" {
		log.Println("running as a client...")
		client.NewClient(int32(*startNumber), int32(*returnLength), *testTimes, *port, *host).Run()
	} else {
		log.Println("running as a server...")
		server.NewServer(*interval, *port, *host).Run()
	}
}
