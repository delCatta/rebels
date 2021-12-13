package service

import (
	"context"
	"fmt"

	"github.com/delCatta/rebels/pb"
)

type Droid struct {
	comms pb.LightSpeedCommsClient
	data  *DroidData
}

func NewDroid(comms pb.LightSpeedCommsClient) *Droid {
	return &Droid{comms: comms, data: NewDroidData()}
}

func (droid *Droid) ToBroker(req *pb.LeiaReq) (*pb.BrokerAmountRes, error) {
	res, err := droid.comms.HowManyRebelsBroker(context.Background(), req)
	if err != nil {
		return nil, err
	}
	fmt.Println(res.GetAmount())
	droid.data.Save(req, res)
	return res, nil
}
