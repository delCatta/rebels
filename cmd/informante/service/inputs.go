package service

import (
	"fmt"
	"strconv"

	"github.com/delCatta/rebels/pb"
)

func GetCommands(droid *Droid) {
	// Given a valid command and its arguments, it sends a request
	// and asks in a loop for another command.
	var command, planet, city string
	var val string
	fmt.Print("Type your command (CNTR+C to exit): ")
	fmt.Scanln(&command, &planet, &city, &val)
	valid, key := IsValidCommand(command)
	if !valid {
		fmt.Printf("Invalid command: %s\n", command)
		GetCommands(droid)
		return
	}
	if key == pb.InformanteReq_NAME_UPDATE {
		request := &pb.InformanteReq{
			Comando:       key,
			NombrePlaneta: planet,
			NombreCiudad:  city,
			NuevoValor: &pb.InformanteReq_NuevaCiudad{
				NuevaCiudad: val,
			},
		}
		err := sendRequest(droid, request)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err.Error())
			GetCommands(droid)
			return
		}

	} else if key == pb.InformanteReq_DELETE {
		request := &pb.InformanteReq{
			Comando:       key,
			NombrePlaneta: planet,
			NombreCiudad:  city,
		}
		err := sendRequest(droid, request)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err.Error())
			GetCommands(droid)
			return
		}

	} else {

		var number uint64
		var err error

		number, err = strconv.ParseUint(val, 10, 64)
		if key == pb.InformanteReq_ADD && err != nil {
			number = 0
			fmt.Printf("Assuming value for %s is: %v\n", command, number)
		} else if err != nil {
			fmt.Printf("Invalid value for the command: %s - Value: %s\n", command, val)
			GetCommands(droid)
			return
		}
		request := &pb.InformanteReq{
			Comando:       key,
			NombrePlaneta: planet,
			NombreCiudad:  city,
			NuevoValor:    &pb.InformanteReq_NuevosRebeldes{NuevosRebeldes: number},
		}
		err = sendRequest(droid, request)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err.Error())
			// Loops
			GetCommands(droid)
			return
		}

	}
	// Loops
	GetCommands(droid)

}
func sendRequest(droid *Droid, request *pb.InformanteReq) error {
	_, err := droid.ToBroker(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func IsValidCommand(command string) (bool, pb.InformanteReq_Command) {
	// Checks if the given command is valid and helps setting it's key.
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
	return false, pb.InformanteReq_ADD // Placeholder.
}
