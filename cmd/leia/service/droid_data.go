package service

import (
	"fmt"
	"sync"

	"github.com/delCatta/rebels/pb"
)

type DroidData struct {
	mutex       *sync.RWMutex
	registers   []*RV
	lastAddress map[string]*pb.FulcrumAddress
}

func NewDroidData() *DroidData {
	return &DroidData{
		mutex:       &sync.RWMutex{},
		registers:   []*RV{},
		lastAddress: make(map[string]*pb.FulcrumAddress),
	}
}

type RV struct {
	req    *pb.LeiaReq
	vector *pb.VectorClock
}

func (data *DroidData) Save(req *pb.LeiaReq, res *pb.BrokerAmountRes) {
	rv := &RV{req: req, vector: res.GetVector()}

	data.mutex.Lock()
	defer data.mutex.Unlock()

	data.registers = append(data.registers, rv)
	data.lastAddress[rv.identifier()] = res.GetAddress()
	// TODO: Notify Saved Data!
	fmt.Println(data.registers)
}

func (rv *RV) identifier() string {
	return rv.req.NombrePlaneta + "-" + rv.req.NombreCiudad
}
