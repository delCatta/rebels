package main

import (
	"log"

	"github.com/delCatta/rebels/cmd/informante/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {

	log.Println("Iniciando Informante...")
	// Conection with Broker
	client := brokerClient("10.6.43.141:3033")
	if client == nil {
		log.Println("Broker not available (Connection Refused)...")
		return
	}
	// Droid helps run the Broker Requests
	droid := service.NewDroid(client)
	// Service starts commands loop.
	service.GetCommands(droid)
}

func brokerClient(ipAddress string) pb.LightSpeedCommsClient {
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	return pb.NewLightSpeedCommsClient(conn)
}
