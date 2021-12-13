package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/delCatta/rebels/pb"
	"google.golang.org/grpc"
)

type FulcrumServer struct {
	reloj    pb.VectorClock
	ftype    string
	planetas map[string]*pb.VectorClock
	fulcrum1 pb.PropagacionCambiosClient
	fulcrum2 pb.PropagacionCambiosClient
	pb.UnimplementedLightSpeedCommsServer
	pb.UnimplementedPropagacionCambiosServer
}

const (
	archivo_log = "registro.log"
)

func NewFulcrumServer(fulcrumType string) *FulcrumServer {
    server := &FulcrumServer{
        reloj:    pb.VectorClock{X: 0, Y: 0, Z: 0},
		ftype:    fulcrumType,
		planetas: make(map[string]*pb.VectorClock),
	}
    if fulcrumType == "Z" {
		server.fulcrum1 = NewFulcrumClient("10.6.43.142")
		server.fulcrum2 = NewFulcrumClient("10.6.43.144")
    }
    return server
}
func NewFulcrumClient(address string) pb.PropagacionCambiosClient {
	conn, err := grpc.Dial(address+":3005", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connecting to %s: %e\n", address, err)
		return nil
	}
	client := pb.NewPropagacionCambiosClient(conn)
	return client
}

func (server *FulcrumServer) InformarFulcrum(ctx context.Context, req *pb.InformanteReq) (*pb.FulcrumRes, error) {
	log_registro, err := os.OpenFile(archivo_log, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		// esto no debiese ocurrir pero en caso que si no hay mucho que se pueda hacer
		fmt.Printf("no se pudo abrir el log de registro: %v", err)
		return nil, err
	}
	defer log_registro.Close()

	switch req.Comando {
	case pb.InformanteReq_ADD:
        fmt.Fprintf(log_registro, "AddCity %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
        fmt.Printf("AddCity %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
		err = server.agregarCiudad(req)

	case pb.InformanteReq_NAME_UPDATE:
		fmt.Fprintf(log_registro, "UpdateName %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevaCiudad())
		fmt.Printf("UpdateName %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevaCiudad())
		err = server.cambiarNombre(req)

	case pb.InformanteReq_NUMBER_UPDATE:
		fmt.Fprintf(log_registro, "UpdateNumber %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
		fmt.Printf("UpdateNumber %v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
		err = server.cambiarValor(req)

	case pb.InformanteReq_DELETE:
		fmt.Fprintf(log_registro, "DeleteCity %v %v\n", req.NombrePlaneta, req.NombreCiudad)
		fmt.Printf("DeleteCity %v %v\n", req.NombrePlaneta, req.NombreCiudad)
		err = server.borrarCiudad(req)
	}

	if err != nil {
		return nil, err
	}
	reloj, _ := server.planetas[req.NombrePlaneta]
	return &pb.FulcrumRes{Vector: reloj}, nil
}

// Esto es porque me di cuenta que el agregar una ciudad necesita verificar que la ciudad no exista
// previamente y eso es básicamente lo que hace el cambiar el valor
func (server *FulcrumServer) agregarCiudad(req *pb.InformanteReq) error {
	return server.cambiarValor(req)
}

// NOTE(lucas): copiado de https://stackoverflow.com/questions/26152901/replace-a-line-in-text-file-golang
func (server *FulcrumServer) cambiarNombre(req *pb.InformanteReq) error {
	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	err := server.sumarComponente(1, existe, req.NombrePlaneta)
	if err != nil {
		return nil
	}

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		// puede darse la situacion donde el planeta tenga un reloj asociado pero no un archivo,
		// esto no debe ser tratado como un error, solo que ninguna operación que genere un archivo
		// ha sido realizada por los informantes sobre este planeta
		if os.IsNotExist(err) {
			return nil
		}
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
			break
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
	_, existe := server.planetas[req.NombrePlaneta]
	err := server.sumarComponente(1, existe, req.NombrePlaneta)

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		// misma situación que en cambiarNombre, si no existe el archivo, no es un error, solo hay
		// que agregar la entrada
		if os.IsNotExist(err) {
			registro_planetario, err := os.OpenFile(req.NombrePlaneta, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer registro_planetario.Close()
			fmt.Fprintf(registro_planetario, "%v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
			return nil
		}
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}

	registro_planetario := string(_registro_planetario)
	entradas := strings.Split(registro_planetario, "\n")

	planeta := ""
	ciudad := ""
	var rebeldes uint64 = 0

	tiene_ciudad := false
	for i, entrada := range entradas {
        if i == len(entradas) - 1 {
            break;
        }
		fmt.Sscanf(entrada, "%v %v %v", &planeta, &ciudad, &rebeldes)
		if ciudad == req.NombreCiudad {
			entradas[i] = fmt.Sprintf("%v %v %v", planeta, req.NombreCiudad, req.GetNuevosRebeldes())
			tiene_ciudad = true
            break;
		}
	}

	if !tiene_ciudad {
		entradas[len(entradas) - 1] = fmt.Sprintf("%v %v %v\n", req.NombrePlaneta, req.NombreCiudad, req.GetNuevosRebeldes())
	}

	out_registro := strings.Join(entradas, "\n")
	err = ioutil.WriteFile(req.NombrePlaneta, []byte(out_registro), 0644)
	if err != nil {
		fmt.Printf("no se pudo registrar el cambio de valor: %v", err)
		return err
	}

	return nil
}

func (server *FulcrumServer) sumarComponente(amount int, existe bool, nombrePlaneta string) error {
	if !existe {
		if server.ftype == "X" {
			server.reloj.X += 1
			server.planetas[nombrePlaneta] = &pb.VectorClock{
				X: server.reloj.X,
				Y: server.reloj.Y,
				Z: server.reloj.Z,
			}
		}
		if server.ftype == "Y" {
			server.reloj.Y += 1
			server.planetas[nombrePlaneta] = &pb.VectorClock{
				X: server.reloj.X,
				Y: server.reloj.Y,
				Z: server.reloj.Z,
			}
		}
		if server.ftype == "Z" {
			server.reloj.Z += 1
			server.planetas[nombrePlaneta] = &pb.VectorClock{
				X: server.reloj.X,
				Y: server.reloj.Y,
				Z: server.reloj.Z,
			}

		}
		return fmt.Errorf("No existe.")
	}

	if server.ftype == "X" {
		server.reloj.X += 1
		server.planetas[nombrePlaneta].X += 1
	}
	if server.ftype == "Y" {
		server.reloj.Y += 1
		server.planetas[nombrePlaneta].Y += 1
	}
	if server.ftype == "Z" {
		server.reloj.Z += 1
		server.planetas[nombrePlaneta].Z += 1
	}
	return nil
}

func (server *FulcrumServer) borrarCiudad(req *pb.InformanteReq) error {
	// revisamos existe ya el registro para el planeta correspondiente
	_, existe := server.planetas[req.NombrePlaneta]
	err := server.sumarComponente(1, existe, req.NombrePlaneta)
	if err != nil {
		return nil
	}

	_registro_planetario, err := ioutil.ReadFile(req.NombrePlaneta)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		fmt.Printf("no se pudo abrir el registro planetario %v\n", err)
		return err
	}

	registro_planetario := string(_registro_planetario)
	entradas := strings.Split(registro_planetario, "\n")

	planeta := ""
	ciudad := ""
	var rebeldes uint64 = 0

    nuevo_registro := []string{}

	for _, entrada := range entradas {
		fmt.Sscanf(entrada, "%v %v %v", &planeta, &ciudad, &rebeldes)
		if ciudad != req.NombreCiudad {
            nuevo_registro = append(nuevo_registro, entrada)
		}
	}

	out_registro := strings.Join(nuevo_registro, "\n")
	err = ioutil.WriteFile(req.NombrePlaneta, []byte(out_registro), 0644)
	if err != nil {
		fmt.Printf("no se pudo registrar el cambio de nombre: %v", err)
		return err
	}

	return nil
}

func (server *FulcrumServer) HowManyRebelsBroker(ctx context.Context, req *pb.LeiaReq) (*pb.BrokerAmountRes, error) {
	registro, err := os.OpenFile(req.NombrePlaneta, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no hay registro de este planeta en este nodo")
		} else {
			return nil, err
		}
	}
	defer registro.Close()

	for {
		var planeta string
		var ciudad string
		var rebeldes uint64
		_, err = fmt.Fscanf(registro, "%v %v %v\n", &planeta, &ciudad, &rebeldes)
		if err != nil {
			break
		}

		if ciudad == req.NombreCiudad {
			return &pb.BrokerAmountRes{
				Vector: server.planetas[req.NombrePlaneta],
				Amount: rebeldes,
			}, nil
		}
	}

	return nil, fmt.Errorf("no hay registro de esta ciudad en este nodo")
}
