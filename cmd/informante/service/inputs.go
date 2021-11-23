package service

import (
	"fmt"

	"github.com/delCatta/rebels/pb"
)

func GetCommands(droid *Droid) {
	var command, planet, city string
	var val uint64
	fmt.Print("Type your command (CNTR+C to exit): ")
	fmt.Scanln(&command, &planet, &city, &val)
	valid, key := IsValidCommand(command)
	if !valid {
		fmt.Printf("Invalid command: %s\n", command)
		GetCommands(droid)
		return
	}
	request := &pb.InformanteReq{
		Comando:       key,
		NombrePlaneta: planet,
		NombreCiudad:  city,
		NuevoValor:    val,
	}
	_, err := droid.ToBroker(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	GetCommands(droid)

}

func IsValidCommand(command string) (bool, pb.InformanteReq_Command) {
	switch command {
	case "AddCity":
		return true, pb.InformanteReq_ADD
	case "UpdateName":
		return true, pb.InformanteReq_NAME_UPDATE
	case "UpdateNumber":
		return true, pb.InformanteReq_NUMBER_UPDATE
	case "DeleteCity":
		return true, pb.InformanteReq_DELETE
	}
	return false, pb.InformanteReq_ADD
}
