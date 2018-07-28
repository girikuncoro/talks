package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "gojek.com/go-points/gopoints"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9001"
)

type goPointsServer struct{}

// GetPoint implements gopoints.GetPoint
func (s *goPointsServer) GetPoint(ctx context.Context, in *pb.PointRequest) (*pb.PointReply, error) {
	point := strconv.Itoa(rand.Intn(100))
	return &pb.PointReply{Point: "Hi " + in.Name + ", you received " + point + " points from GO-JEK"}, nil
}

func main() {
	log.Printf("GO-POINTS server is listening on port %s ...", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGoPointsServer(grpcServer, &goPointsServer{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
