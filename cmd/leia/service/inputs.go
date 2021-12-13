package service

import (
	"fmt"

	"github.com/delCatta/rebels/pb"
)

func GetCommands(droid *Droid) {
	var command, planet, city string
	fmt.Print("Type your command (CNTR+C to exit): ")
	fmt.Scanln(&command, &planet, &city)
	valid := IsValidCommand(command)
	if !valid {
		fmt.Printf("Invalid command: %s\n", command)
		GetCommands(droid)
		return
	}
	request := &pb.LeiaReq{
		NombrePlaneta: planet,
		NombreCiudad:  city,
	}
	_, err := droid.ToBroker(request)
	if err != nil {
		fmt.Println(err)
	}
	GetCommands(droid)
}

func IsValidCommand(command string) bool {
	switch command {
	case "GetNumberRebelds":
		return true
	}
	return false
}
