package service

import (
	"fmt"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type FulcrumServer struct {
	fulcrum1 *pb.LightSpeedCommsClient
	fulcrum2 *pb.LightSpeedCommsClient
	fulcrum3 *pb.LightSpeedCommsClient
	pb.UnimplementedLightSpeedCommsServer
}

func NewFulcrumServer() *FulcrumServer {
	// TODO: No conectarse consigo mismo...
	return &FulcrumServer{
		fulcrum1: NewFulcrumClient("IP 1 SIN PUERTO"),
		fulcrum2: NewFulcrumClient("IP 2 SIN PUERTO"),
		fulcrum3: NewFulcrumClient("IP 3 SIN PUERTO"),
	}
}
func NewFulcrumClient(address string) *pb.LightSpeedCommsClient {
	conn, err := grpc.Dial(address+":3005", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connecting to %s: %e\n", address, err)
		return nil
	}
	client := pb.NewLightSpeedCommsClient(conn)
	return &client
}

// TODO: Proto Comunications con Broker para redireccionar

// TODO: Proto Comunications con otro Fulcrum para merges
