package main

import (
	"log"
	"os"
	"time"

	pb "gojek.com/go-points/gopoints"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:9001"
	defaultName = "Nadiem"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGoPointsClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetPoint(ctx, &pb.PointRequest{Name: name})
	if err != nil {
		log.Fatalf("could not receive point: %v", err)
	}
	log.Printf("GO-POINTS: %s", resp.Point)
}
