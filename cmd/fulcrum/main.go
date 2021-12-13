package main

import (
	"log"
	"net"
    // "sync"

	"github.com/delCatta/rebels/cmd/fulcrum/service"
	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp4", "localhost:3005") // Puerto 3005!
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fulcrumServer := service.NewFulcrumServer() 
	pb.RegisterLightSpeedCommsServer(grpcServer, fulcrumServer)

    // 2 de los fulcrum necesitan esta linea
    pb.RegisterPropagacionCambiosServer(grpcServer, fulcrumServer)

    // el tercero necesita esto:
    // wg := new(sync.WaitGroup)
    // wg.Add(1)

    // go func() {
    //     fulcrumServer.propagarCambios()
    //     wg.Done()
    // }()
	grpcServer.Serve(lis)

    // wg.Wait()
}
