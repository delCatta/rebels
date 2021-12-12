package service

import (
	"context"
	"fmt"
	"os"
    "time"

	"github.com/delCatta/rebels/pb"
)

type RegistroPlanetario struct {
	ciudades map[string][]*pb.InformanteReq
}

func (server *FulcrumServer) propagarDatos() {
    for {
        time.Sleep(2 * time.Minute)
        log1, err1 := server.getLog()
        log2, err2 := server.fulcrum1.MergeBegin(context.Background(), &pb.MergeBeginReq{})
        log3, err3 := server.fulcrum2.MergeBegin(context.Background(), &pb.MergeBeginReq{})
        if err1 != nil || err2 != nil || err3 != nil {
            fmt.Printf("no se ha podido obtener un log\n");
            continue
        }

        merge_end, err := merge([3]*pb.MergeBeginRes{log1, log2, log3})
        if err != nil {
            fmt.Printf("no ha sido posible completar el merge: %v\n", err)
            continue
        }

        server.integrateChanges(merge_end)
        server.fulcrum1.MergeEnd(context.Background(), merge_end)
        server.fulcrum2.MergeEnd(context.Background(), merge_end)
    }
}

func merge(logs [3]*pb.MergeBeginRes) (*pb.MergeEndReq, error) {
	merge_result := &pb.MergeEndReq{}

	_merge_result := make(map[string]RegistroPlanetario)

	for _, log := range logs {
		for _, comando := range log.Changelog {
			if _, existe := _merge_result[comando.NombrePlaneta]; !existe {
				_merge_result[comando.NombrePlaneta] = RegistroPlanetario{
					ciudades: make(map[string][]*pb.InformanteReq),
				}
			}
			
			// NOTE(lucas): esto es para no tener una linea de 200 caracteres
			ciudad :=_merge_result[comando.NombrePlaneta].ciudades[comando.NombreCiudad] 

			if len(ciudad) != 0 {
				ultimo_comando := ciudad[len(ciudad) - 1].Comando 

				// NOTE/TODO(lucas): este criterio es arbitrario y bastante propenso a errores
				// pero ningun criterio que se me ocurre no los tiene asique voy a dejar este por
				// mientras
				if  (ultimo_comando != pb.InformanteReq_DELETE) && (ultimo_comando != pb.InformanteReq_NAME_UPDATE) {
					ciudad = append(ciudad, comando)
					_merge_result[comando.NombrePlaneta].ciudades[comando.NombreCiudad] = ciudad
				}
			}

		}
	}

	return merge_result, nil
}

func (server *FulcrumServer) getLog() (*pb.MergeBeginRes, error) {
    res := &pb.MergeBeginRes{ Reloj: &server.reloj }
    registro_log, err := os.OpenFile(archivo_log, os.O_RDONLY, 0644)
    if err != nil {
        fmt.Printf("no se ha podido abrir el archivo de registro: %v\n", err)
        return nil, err
    }
    defer registro_log.Close()
    for {
        comando := ""
        planeta := ""
        ciudad := ""
        _, err = fmt.Fscanf(registro_log, "%v %v %v ", &comando, &planeta, &ciudad)
        if err != nil  {
            break;
        }
        switch comando {
        case "AddCity":
            var rebeldes uint64            
            fmt.Fscanf(registro_log, "%v\n", &rebeldes)
            res.Changelog = append(res.Changelog, &pb.InformanteReq{
                Comando: pb.InformanteReq_ADD,
                NombrePlaneta: planeta,
                NombreCiudad: ciudad,
                NuevoValor: &pb.InformanteReq_NuevosRebeldes{NuevosRebeldes: rebeldes},
            })
            break;
        case "DeleteCity":
            fmt.Fscanf(registro_log, "\n")
            res.Changelog = append(res.Changelog, &pb.InformanteReq{
                Comando: pb.InformanteReq_DELETE,
                NombrePlaneta: planeta,
                NombreCiudad: ciudad,
            })
            break;
        case "UpdateValue":
            var rebeldes uint64            
            fmt.Fscanf(registro_log, "%v\n", &rebeldes)
            res.Changelog = append(res.Changelog, &pb.InformanteReq{
                Comando: pb.InformanteReq_NUMBER_UPDATE,
                NombrePlaneta: planeta,
                NombreCiudad: ciudad,
                NuevoValor: &pb.InformanteReq_NuevosRebeldes{NuevosRebeldes: rebeldes},
            })
            break;
        case "UpdateName":
            var nueva_ciudad string            
            fmt.Fscanf(registro_log, "%v\n", &nueva_ciudad)
            res.Changelog = append(res.Changelog, &pb.InformanteReq{
                Comando: pb.InformanteReq_NAME_UPDATE,
                NombrePlaneta: planeta,
                NombreCiudad: ciudad,
                NuevoValor: &pb.InformanteReq_NuevaCiudad{NuevaCiudad: nueva_ciudad},
            })
            break;
        default:
            fmt.Printf("comando desconocido: %v\n", comando)
            break;
        }
    }
	return res, nil
}

func (server *FulcrumServer) integrateChanges(req *pb.MergeEndReq) error {
    for _, comando := range req.Changelog {
        var err error
        switch comando.Comando {
        case pb.InformanteReq_ADD:
            err = server.agregarCiudad(comando)
            break;
        case pb.InformanteReq_NAME_UPDATE:
            err = server.cambiarNombre(comando)
            break;
        case pb.InformanteReq_NUMBER_UPDATE:
             err = server.cambiarValor(comando)
            break;
        case pb.InformanteReq_DELETE:
            err = server.borrarCiudad(comando)
            break;
        }

        if err != nil {
            fmt.Println(err)
            return err
        }
    }
    os.Remove(archivo_log)
    return nil
}

func (server *FulcrumServer) MergeBegin(ctx context.Context, req *pb.MergeBeginReq) (*pb.MergeBeginRes, error) {
    return server.getLog()
}

func (server *FulcrumServer) MergeEnd(ctx context.Context, req *pb.MergeEndReq) (*pb.MergeEndRes, error) {
    err := server.integrateChanges(req)
    if err != nil {
        return nil, err
    }
	return &pb.MergeEndRes{}, nil
}

func max3(a int32, b int32, c int32) int32 {
	return max(a, max(b,c))
}

// NOTE(lucas): go no tiene un operador ternario asique esto tiene que ser hecho de la manera fea
func max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
