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
	grpcServer := grpc.NewServer()                             // Crea un servidor de grpc
	brokerServer := service.NewBrokerServer()                  // Crea un servidor de broker
	pb.RegisterLightSpeedCommsServer(grpcServer, brokerServer) // Registra el servidor de grpc en el servidor de broker
	grpcServer.Serve(lis)                                      // Inicia el servidor de grpc
}
