package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type BrokerServer struct { // Servidor de broker con sus respectivos clientes de fulcrum
	fulcrum1  *pb.LightSpeedCommsClient
	fulcrum2  *pb.LightSpeedCommsClient
	fulcrum3  *pb.LightSpeedCommsClient
	addresses map[*pb.LightSpeedCommsClient]*pb.FulcrumAddress
	pb.UnimplementedLightSpeedCommsServer
}

func NewBrokerServer() *BrokerServer { // Funcion para crear un servidor de broker
	brokerServer := &BrokerServer{
		addresses: map[*pb.LightSpeedCommsClient]*pb.FulcrumAddress{},
	}
	brokerServer.fulcrum1 = NewFulcrumClient("10.6.43.142", brokerServer) // Asocia un cliente de fulcrum1 con la ip de la maquina 2
	brokerServer.fulcrum2 = NewFulcrumClient("10.6.43.143", brokerServer) // Asocia un cliente de fulcrum1 con la ip de la maquina 3
	brokerServer.fulcrum3 = NewFulcrumClient("10.6.43.144", brokerServer) // Asocia un cliente de fulcrum1 con la ip de la maquina 4
	return brokerServer
}
func NewFulcrumClient(address string, server *BrokerServer) *pb.LightSpeedCommsClient { // Funcion para asociar un cliente de fulcrum
	conn, err := grpc.Dial(address+":3005", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connecting to %s: %e\n", address, err)
		return nil
	}
	client := pb.NewLightSpeedCommsClient(conn)                      // Se crea un cliente de fulcrum
	server.addresses[&client] = &pb.FulcrumAddress{Address: address} // Se asocia la direccion del cliente de fulcrum con el servidor
	return &client
}

func (server *BrokerServer) InformarBroker(ctx context.Context, req *pb.InformanteReq) (*pb.BrokerRes, error) { // Funcion para informar a un cliente de fulcrum
	fulcrumClient := server.pickAClient()      // Selecciona un cliente de fulcrum aleatoriamente
	address := server.addresses[fulcrumClient] // Obtiene la direccion del cliente de fulcrum
	response := &pb.BrokerRes{                 // Crea una respuesta
		Address: address,
	}
	return response, nil

}

func (server *BrokerServer) HowManyRebelsBroker(ctx context.Context, req *pb.LeiaReq) (*pb.BrokerAmountRes, error) { // Funcion para informar la cantiad de rebeldes proveniente de la request de leia
	for i := 0; i < 3; i++ { // Se itera por cada cliente de fulcrum hasta encontrar el correcto
		var fulcrumClient *pb.LightSpeedCommsClient
		if i == 0 {
			fmt.Println("fulcrum1")
			fulcrumClient = server.fulcrum1
		}
		if i == 1 {
			fmt.Println("fulcrum2")
			fulcrumClient = server.fulcrum2
		}
		if i == 2 {
			fmt.Println("fulcrum3")
			fulcrumClient = server.fulcrum3
		}

		res, err := (*fulcrumClient).HowManyRebelsBroker(ctx, req) // Se llama a la funcion de HowManyRebelsBroker del cliente de fulcrum con la request de leia
		if err != nil {
			fmt.Printf("Error calling HowManyRebels: %v\n", err.Error())
			if i == 2 {
				return nil, err
			}
			continue
		}
		response := &pb.BrokerAmountRes{ // Se crea una respuesta
			Address: server.addresses[fulcrumClient], // Se asocia la direccion del cliente de fulcrum
			Vector:  res.GetVector(),                 // Se asocia el vector de la respuesta
			Amount:  res.GetAmount(),                 // Se asocia la cantidad de rebeldes
		}
		return response, nil
	}
	return nil, nil
}

func (server *BrokerServer) pickAClient() *pb.LightSpeedCommsClient { // Funcion para seleccionar un cliente de fulcrum aleatoriamente
	rand.Seed(time.Now().UnixNano()) // Se inicializa la semilla
	random_int := rand.Intn(3)       // Se genera un numero aleatorio entre 0 y 2
	if random_int == 0 {             // Si el numero aleatorio es 0 se escoge el fulcrum 1
		if server.fulcrum1 != nil {
			return server.fulcrum1
		}
	}
	if random_int == 1 { // Si el numero aleatorio es 1 se escoge el fulcrum 2
		if server.fulcrum2 != nil {
			return server.fulcrum2
		}
	}
	if random_int == 2 { // Si el numero aleatorio es 2 se escoge el fulcrum 3
		if server.fulcrum3 != nil {
			return server.fulcrum3
		}
	}
	return nil
}
