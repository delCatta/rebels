package main

import (
	"log"
	"net"

	"github.com/delCatta/rebels/cmd/fulcrum/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:3005") // Puerto 3005!
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fulcrumServer := service.NewFulcrumServer() // TODO: Implement this FulcrumServer
	pb.RegisterLightSpeedCommsServer(grpcServer, fulcrumServer)
	grpcServer.Serve(lis)
}
