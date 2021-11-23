package service

import (
	"context"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type Droid struct {
	comms pb.LightSpeedCommsClient
	data  *DroidData
}

func NewDroid(comms pb.LightSpeedCommsClient) *Droid {
	return &Droid{comms: comms, data: NewDroidData()}
}

func (droid *Droid) ToBroker(req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	res, err := droid.comms.InformarBroker(context.Background(), &pb.InformanteReq{
		Comando: pb.InformanteReq_ADD,
	})
	if err != nil {
		return nil, err
	}
	return droid.toServer(res.Address, req)
}
func (droid *Droid) toServer(address *pb.FulcrumAddress, req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	fulcrum_client, err := fulcrumClient(address)
	if err != nil {
		return nil, err
	}
	res, err := fulcrum_client.InformarFulcrum(context.Background(), &pb.InformanteReq{
		Comando: pb.InformanteReq_ADD,
	})
	if err != nil {
		return nil, err
	}
	// TODO: Save Received Information.
	return res, nil
}

func fulcrumClient(ipAddress *pb.FulcrumAddress) (pb.LightSpeedCommsClient, error) {
	conn, err := grpc.Dial(ipAddress.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewLightSpeedCommsClient(conn), nil
}
