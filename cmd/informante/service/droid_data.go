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
	req    *pb.InformanteReq
	vector *pb.VectorClock
}

func (data *DroidData) Save(req *pb.InformanteReq, res *pb.FulcrumRes, address *pb.FulcrumAddress) {
	rv := &RV{req: req, vector: res.GetVector()}
	fmt.Println(res, res.GetVector())
	data.mutex.Lock()
	defer data.mutex.Unlock()

	data.registers = append(data.registers, rv)
	data.lastAddress[rv.identifier()] = address
	fmt.Println("Request received successfully.")
}

func (rv *RV) identifier() string {
	return rv.req.NombrePlaneta + "-" + rv.req.NombreCiudad
}
