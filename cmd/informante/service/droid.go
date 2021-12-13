package service

import (
	"context"
	"fmt"

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
	res, err := droid.comms.InformarBroker(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return droid.toServer(res.GetAddress(), req)
}
func (droid *Droid) toServer(address *pb.FulcrumAddress, req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	fulcrum_client, err := fulcrumClient(address)
	if err != nil {
		return nil, err
	}
	res, err := fulcrum_client.InformarFulcrum(context.Background(), req)
	if err != nil {
		return nil, err
	}
	droid.data.Save(req, res, address)
	return res, nil
}

func fulcrumClient(ipAddress *pb.FulcrumAddress) (pb.LightSpeedCommsClient, error) {
	fmt.Printf("Connecting with Fulcrum: %s \n", ipAddress.GetAddress())
	conn, err := grpc.Dial(ipAddress.GetAddress()+":3005", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewLightSpeedCommsClient(conn), nil
}
