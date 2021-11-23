package main

import (
	"log"

	"github.com/delCatta/rebels/cmd/informante/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {
	// Connecting to Broker
	log.Println("Iniciando Informante...")
	client := brokerClient("localhost:3033")
	if client == nil {
		log.Println("Broker not available (Connection Refused)...")
		return
	}
	droid := service.NewDroid(client)
	service.GetCommands(droid)
}

func brokerClient(ipAddress string) pb.LightSpeedCommsClient {
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	return pb.NewLightSpeedCommsClient(conn)
}
