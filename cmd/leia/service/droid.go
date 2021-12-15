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
	// Creates a new Droid Structure that stores data.

	return &Droid{comms: comms, data: NewDroidData()}
}

func (droid *Droid) ToBroker(req *pb.LeiaReq) (*pb.BrokerAmountRes, error) {
	// Informs Broker with a Request
	res, err := droid.comms.HowManyRebelsBroker(context.Background(), req)
	if err != nil {
		return nil, err
	}
	fmt.Println(res.GetAmount())
	// Stores the data in the droid's memory.
	droid.data.Save(req, res)
	return res, nil
}
