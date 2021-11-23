package service

import (
	"context"
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

func (server *FulcrumServer) InformarFulcrum(ctx context.Context, req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	// TODO: Almacenar la request y enviar el vector en la respuesta.
	response := &pb.FulcrumRes{
		Vector: nil,
	}
	return response, nil
}

// TODO: Proto Comunications con Broker para redireccionar

// TODO: Proto Comunications con otro Fulcrum para merges
