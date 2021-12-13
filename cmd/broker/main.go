package main

import (
	"log"
	"net"

	"github.com/delCatta/rebels/cmd/broker/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:3033") // Puerto 3003!
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	brokerServer := service.NewBrokerServer() // TODO: Implement this BrokerServer
	pb.RegisterLightSpeedCommsServer(grpcServer, brokerServer)
	grpcServer.Serve(lis)
}
