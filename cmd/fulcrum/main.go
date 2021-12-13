package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	// "sync"

	"github.com/delCatta/rebels/cmd/fulcrum/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Se necesita identificar el tipo de Fulcrum.")
	}
	fulcrumType := os.Args[1]

	lis, err := net.Listen("tcp4", ":3005") // Puerto 3005!
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fulcrumServer := service.NewFulcrumServer(fulcrumType)
	pb.RegisterLightSpeedCommsServer(grpcServer, fulcrumServer)

	wg := new(sync.WaitGroup)
	if fulcrumType == "X" || fulcrumType == "Y" {
		pb.RegisterPropagacionCambiosServer(grpcServer, fulcrumServer)
	} else {

		wg.Add(1)

		go func() {
			fulcrumServer.PropagarCambios()
			wg.Done()
		}()
	}
	grpcServer.Serve(lis)
	wg.Wait()

}
