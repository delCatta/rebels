package service

import (
	"context"
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type FulcrumServer struct {
	fulcrum1 *pb.LightSpeedCommsClient
	fulcrum2 *pb.LightSpeedCommsClient
	fulcrum3 *pb.LightSpeedCommsClient
	reloj    pb.VectorClock
	planetas map[string]pb.VectorClock
	pb.UnimplementedLightSpeedCommsServer
}

func NewFulcrumServer() *FulcrumServer {
	// TODO: No conectarse consigo mismo...
	return &FulcrumServer{
		fulcrum1: NewFulcrumClient("IP 1 SIN PUERTO"),
		fulcrum2: NewFulcrumClient("IP 2 SIN PUERTO"),
		fulcrum3: NewFulcrumClient("IP 3 SIN PUERTO"),
		reloj:    pb.VectorClock{},
		planetas: make(map[string]pb.VectorClock),
	}
}
func NewFulcrumClient(address string) *pb.LightSpeedCommsClient {
	conn, err := grpc.Dial(address+":3005", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connecting to %s: %e\n", address, err)
		return nil
	}
	client := pb.NewLightSpeedCommsClient(conn)
	return &client
}

func (server *FulcrumServer) InformarFulcrum(ctx context.Context, req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	// TODO: Almacenar la request y enviar el vector en la respuesta.

	log_registro, err := os.OpenFile("registro.log", os.O_APPEND |os.O_CREATE, 0644)
	if err != nil {
		// esto no debiese ocurrir pero en caso que si no hay mucho que se pueda hacer
		fmt.Printf("no se pudo abrir el log de registro: %v", err)
		return nil, err
	}
	defer log_registro.Close()

	switch req.Comando {
	case pb.InformanteReq_ADD:
		// TODO(lucas): ¿logear aun si falla registrar la request?, de momento si
		fmt.Fprintf(log_registro, "AddCity %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
		err = server.agregarCiudad(req)
		break;

	case pb.InformanteReq_NAME_UPDATE:
		fmt.Fprintf(log_registro, "UpdateName %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevaCiudad())
		err = server.cambiarNombre(req)
		break;

	case pb.InformanteReq_NUMBER_UPDATE:
		fmt.Fprintf(log_registro, "UpdateNumber %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
		err = server.cambiarValor(req)
		break;

	case pb.InformanteReq_DELETE:
		fmt.Fprintf(log_registro, "DeleteCity %v %v\n", req.NombrePlaneta, req.NombreCiudad)
		err = server.borrarCiudad(req)
		break;
	}

	if err != nil {
		return nil, err
	}
	reloj, _ := server.planetas[req.NombrePlaneta]
	return &pb.FulcrumRes{ Vector: &reloj }, nil
}


func (server *FulcrumServer) agregarCiudad(req *pb.InformanteReq) error {
	registro_planetario, err := os.OpenFile(req.NombrePlaneta, os.O_APPEND | os.O_CREATE, 0644);
	if err != nil {
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}
	defer registro_planetario.Close()

	fmt.Fprintf(registro_planetario, "%v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.NuevoValor)

	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	if !existe {
		// TODO(lucas): sumarle 1 a la componente correspondiente al servidor
		server.reloj.X += 1
		server.planetas[req.NombrePlaneta] = pb.VectorClock{
			X: server.reloj.X + 1,
			Y: server.reloj.Y,
			Z: server.reloj.Z,
		}
	}
	return nil
}

// NOTE(lucas): copiado de https://stackoverflow.com/questions/26152901/replace-a-line-in-text-file-golang
func (server *FulcrumServer) cambiarNombre(req *pb.InformanteReq) error {
	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	if !existe {
		// TODO(lucas): sumarle 1 a la componente correspondiente al servidor
		server.reloj.X += 1
		server.planetas[req.NombrePlaneta] = pb.VectorClock{
			X: server.reloj.X + 1,
			Y: server.reloj.Y,
			Z: server.reloj.Z,
		}

		// en este caso si no habia registro del planeta entonces no hay una ciudad a la que
		// cambiarle el nombre
		return nil
	}

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}

	registro_planetario := string(_registro_planetario)
	entradas := strings.Split(registro_planetario, "\n")


	planeta := ""
	ciudad := ""
	var rebeldes uint64 = 0

	for i, entrada := range entradas {
		_, err := fmt.Sscanf(entrada, "%v %v %v", &planeta, &ciudad, &rebeldes)
		if err != nil {
			break;
		}
		if ciudad == req.NombreCiudad {
			entradas[i] = fmt.Sprintf("%v %v %v", planeta, req.GetNuevaCiudad(), rebeldes)
		}
	}

	out_registro := strings.Join(entradas, "\n")
	err = ioutil.WriteFile(req.NombrePlaneta, []byte(out_registro), 0644)
	if err != nil {
		fmt.Printf("no se pudo registrar el cambio de nombre: %v", err)
		return err
	}

	return nil
}

func (server *FulcrumServer) cambiarValor(req *pb.InformanteReq) error {
	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	if !existe {
		// TODO(lucas): sumarle 1 a la componente correspondiente al servidor
		server.reloj.X += 1
		server.planetas[req.NombrePlaneta] = pb.VectorClock{
			X: server.reloj.X + 1,
			Y: server.reloj.Y,
			Z: server.reloj.Z,
		}

		// en este caso si no habia registro del planeta entonces no hay una ciudad a la que
		// cambiarle el valor
		return nil
	}

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}

	registro_planetario := string(_registro_planetario)
	entradas := strings.Split(registro_planetario, "\n")


	planeta := ""
	ciudad := ""
	var rebeldes uint64 = 0

	for i, entrada := range entradas {
		_, err := fmt.Sscanf(entrada, "%v %v %v", &planeta, &ciudad, &rebeldes)
		if err != nil {
			break;
		}
		if ciudad == req.NombreCiudad {
			entradas[i] = fmt.Sprintf("%v %v %v", planeta, req.NombreCiudad, req.GetNuevosRebeldes())
		}
	}

	out_registro := strings.Join(entradas, "\n")
	err = ioutil.WriteFile(req.NombrePlaneta, []byte(out_registro), 0644)
	if err != nil {
		fmt.Printf("no se pudo registrar el cambio de nombre: %v", err)
		return err
	}

	return nil
}

func (server *FulcrumServer) borrarCiudad(req *pb.InformanteReq) error {
	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	if !existe {
		// TODO(lucas): sumarle 1 a la componente correspondiente al servidor
		server.reloj.X += 1
		server.planetas[req.NombrePlaneta] = pb.VectorClock{
			X: server.reloj.X + 1,
			Y: server.reloj.Y,
			Z: server.reloj.Z,
		}

		// en este caso si no habia registro del planeta entonces no hay una ciudad que borrar
		return nil
	}

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}

	registro_planetario := string(_registro_planetario)
	entradas := strings.Split(registro_planetario, "\n")


	planeta := ""
	ciudad := ""
	var rebeldes uint64 = 0

	for i, entrada := range entradas {
		_, err := fmt.Sscanf(entrada, "%v %v %v", &planeta, &ciudad, &rebeldes)
		if err != nil {
			break;
		}
		if ciudad == req.NombreCiudad {
			entradas[i] = "" // NOTE(lucas): no se si esto resultará en una linea vacia
		}
	}

	out_registro := strings.Join(entradas, "\n")
	err = ioutil.WriteFile(req.NombrePlaneta, []byte(out_registro), 0644)
	if err != nil {
		fmt.Printf("no se pudo registrar el cambio de nombre: %v", err)
		return err
	}

	return nil
}

// TODO: Proto Comunications con Broker para redireccionar

// TODO: Proto Comunications con otro Fulcrum para merges
