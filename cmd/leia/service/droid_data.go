package service

import (
	"sync"

	"github.com/delCatta/rebels/pb"
)

type DroidData struct {
	mutex       *sync.RWMutex
	registers   []*RV
	lastAddress map[string]*pb.FulcrumAddress
}

func NewDroidData() *DroidData {
	// Droid data is where the Droid's data is stored! Acts as memory.
	return &DroidData{
		mutex:       &sync.RWMutex{},
		registers:   []*RV{},
		lastAddress: make(map[string]*pb.FulcrumAddress),
	}
}

type RV struct {
	// Data structured to be stored.
	req    *pb.LeiaReq
	vector *pb.VectorClock
}

// Method to store the RV Data structure
func (data *DroidData) Save(req *pb.LeiaReq, res *pb.BrokerAmountRes) {
	rv := &RV{req: req, vector: res.GetVector()}

	data.mutex.Lock()
	defer data.mutex.Unlock()

	data.registers = append(data.registers, rv)
	data.lastAddress[rv.identifier()] = res.GetAddress()
}

func (rv *RV) identifier() string {
	return rv.req.NombrePlaneta + "-" + rv.req.NombreCiudad
}
