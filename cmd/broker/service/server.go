package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type BrokerServer struct {
	fulcrum1  *pb.LightSpeedCommsClient
	fulcrum2  *pb.LightSpeedCommsClient
	fulcrum3  *pb.LightSpeedCommsClient
	addresses map[*pb.LightSpeedCommsClient]*pb.FulcrumAddress
	pb.UnimplementedLightSpeedCommsServer
}

func NewBrokerServer() *BrokerServer {
	brokerServer := &BrokerServer{
		addresses: map[*pb.LightSpeedCommsClient]*pb.FulcrumAddress{},
	}
	brokerServer.fulcrum1 = NewFulcrumClient("IP 1 SIN PUERTO", brokerServer)
	brokerServer.fulcrum2 = NewFulcrumClient("IP 2 SIN PUERTO", brokerServer)
	brokerServer.fulcrum3 = NewFulcrumClient("IP 3 SIN PUERTO", brokerServer)
	return brokerServer
}
func NewFulcrumClient(address string, server *BrokerServer) *pb.LightSpeedCommsClient {
	conn, err := grpc.Dial(address+":3005", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connecting to %s: %e\n", address, err)
		return nil
	}
	client := pb.NewLightSpeedCommsClient(conn)
	server.addresses[&client] = &pb.FulcrumAddress{Address: address}
	return &client

}

func (server *BrokerServer) InformarBroker(ctx context.Context, req *pb.InformanteReq) (*pb.BrokerRes, error) {
	fulcrumClient := server.pickAClient()
	address := server.addresses[fulcrumClient]
	response := &pb.BrokerRes{
		Address: address,
	}
	return response, nil

}

func (server *BrokerServer) HowManyRebelsBroker(ctx context.Context, req *pb.LeiaReq) (*pb.BrokerAmountRes, error) {
	fulcrumClient := server.pickAClient()
	request := &pb.LeiaReq{}
	res, err := (*fulcrumClient).HowManyRebelsBroker(ctx, request)
	if err != nil {
		fmt.Printf("Error calling HowManyRebels: %e\n", err)
		return nil, err
	}
	response := &pb.BrokerAmountRes{
		Address: server.addresses[fulcrumClient],
		Vector:  res.GetVector(),
		Amount:  res.GetAmount(),
	}
	return response, nil
}

func (server *BrokerServer) pickAClient() *pb.LightSpeedCommsClient {
	rand.Seed(time.Now().UnixNano())
	random_int := rand.Intn(3)
	if random_int == 0 {
		if server.fulcrum1 != nil {
			return server.fulcrum1
		}
	}
	if random_int == 1 {
		if server.fulcrum2 != nil {
			return server.fulcrum2
		}
	}
	if random_int == 2 {
		if server.fulcrum3 != nil {
			return server.fulcrum3
		}
	}
	return nil
}
