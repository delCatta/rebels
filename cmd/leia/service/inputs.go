package service

import (
	"fmt"

	"github.com/delCatta/rebels/pb"
)

func GetCommands(droid *Droid) {
	// Given a valid command and its arguments, it sends a request
	// and asks in a loop for another command.
	var command, planet, city string
	fmt.Print("Type your command (CNTR+C to exit): ")
	fmt.Scanln(&command, &planet, &city)
	valid := IsValidCommand(command)
	if !valid {
		fmt.Printf("Invalid command: %s\n", command)
		// Loop
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
	// Loop
	GetCommands(droid)
}

func IsValidCommand(command string) bool {
	// Checks if the given command is valid and helps setting it's key.
	switch command {
	case "GetNumberRebelds":
		return true
	}
	return false
}
