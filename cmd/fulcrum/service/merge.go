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

func (server *FulcrumServer) PropagarCambios() {
	for {
		time.Sleep(2 * time.Minute)
		log1, err1 := server.obtenerLog()
		log2, err2 := server.fulcrum1.MergeBegin(context.Background(), &pb.MergeBeginReq{})
		log3, err3 := server.fulcrum2.MergeBegin(context.Background(), &pb.MergeBeginReq{})
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Printf("no se ha podido obtener un log\n")
			continue
		}

		merge_end, err := merge([3]*pb.MergeBeginRes{log1, log2, log3})
		if err != nil {
			fmt.Printf("no ha sido posible completar el merge: %v\n", err)
			continue
		}

		merge_end.Reloj = &pb.VectorClock{
			X: max3(log1.Reloj.X, log2.Reloj.X, log3.Reloj.X),
			Y: max3(log1.Reloj.Y, log2.Reloj.Y, log3.Reloj.Y),
			Z: max3(log1.Reloj.Z, log2.Reloj.Z, log3.Reloj.Z),
		}

		server.integrarCambios(merge_end)
		server.fulcrum1.MergeEnd(context.Background(), merge_end)
		server.fulcrum2.MergeEnd(context.Background(), merge_end)
	}
}

func merge(logs [3]*pb.MergeBeginRes) (*pb.MergeEndReq, error) {
	merge_result := &pb.MergeEndReq{}
	tmp := mergeEntre2Logs(logs[0].Changelog, logs[1].Changelog)
	merge_result.Changelog = mergeEntre2Logs(tmp, logs[2].Changelog)
	return merge_result, nil
}

func mergeEntre2Logs(log1 []*pb.InformanteReq, log2 []*pb.InformanteReq) []*pb.InformanteReq {
	len1 := len(log1)
	len2 := len(log2)
	res := make([]*pb.InformanteReq, 0, len1+len2)

	mergeado_hasta := [2]int{0, 0}

	for i, cambio := range log1 {
		switch cambio.Comando {
		case pb.InformanteReq_ADD:
			// si log2 hace mencion de la ciudad a la que se hizo add, hay que mergear todos los
			// cambios anteriores a que eso ocurra
			pos := primeraAparicion(log2[mergeado_hasta[1]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, mergeEntre2Logs(log2[mergeado_hasta[1]:pos], log1[mergeado_hasta[0]:i])...)
				res = append(res, log1[i])
				res = append(res, log2[pos])
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}
			break
		case pb.InformanteReq_DELETE:
			// si log2 hace mencion de la ciudad borrada, hay que encontrar la ultima aparicion de
			// esa ciudad y  mergear los cambios hasta ese punto
			pos := ultimaAparicion(log2[mergeado_hasta[1]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, mergeEntre2Logs(log2[mergeado_hasta[1]:pos], log1[mergeado_hasta[0]:i])...)
				res = append(res, log2[pos])
				res = append(res, log1[i])
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}
			break
		case pb.InformanteReq_NAME_UPDATE:
			mergeado_actual := false

			// cualquier cambio con el nombre viejo debe ser mergeado antes
			pos := ultimaAparicion(log2[mergeado_hasta[1]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, mergeEntre2Logs(log2[mergeado_hasta[1]:pos], log1[mergeado_hasta[0]:i])...)
				mergeado_actual = true

				res = append(res, log2[pos])
				res = append(res, log1[i])
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}

			// los anteriores a que se ocupe el nuevo nombre deben ser mergeados despues
			pos = primeraAparicion(log2[mergeado_hasta[1]:], cambio.GetNuevaCiudad())
			if pos != -1 {
				if !mergeado_actual {
					res = append(res, mergeEntre2Logs(log2[mergeado_hasta[1]:pos], log1[mergeado_hasta[0]:i])...)
				}
				res = append(res, log1[i])
				res = append(res, log2[pos])
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}
			break
		case pb.InformanteReq_NUMBER_UPDATE:
			// este caso es trivial ya que por su cuenta no nos dice nada sobre el orden en que
			// ocurrieron otros eventos
			break
		}
	}

	// mergeamos los cambios restantes.
	// en este punto no hay nada que puede generar conflictos en log1
	for i := mergeado_hasta[1]; i < len2; i++ {
		cambio := log2[i]
		switch cambio.Comando {
		case pb.InformanteReq_ADD:
			pos := primeraAparicion(log1[mergeado_hasta[0]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, log2[mergeado_hasta[1]:i+1]...)
				res = append(res, log1[mergeado_hasta[0]:pos+1]...)
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}
			break
		case pb.InformanteReq_DELETE:
			pos := ultimaAparicion(log1[mergeado_hasta[0]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, log1[mergeado_hasta[0]:pos+1]...)
				res = append(res, log2[mergeado_hasta[1]:i+1]...)
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}
			break
		case pb.InformanteReq_NAME_UPDATE:
			mergeado_actual := false
			pos := ultimaAparicion(log1[mergeado_hasta[0]:], cambio.NombreCiudad)
			if pos != -1 {
				res = append(res, log1[mergeado_hasta[0]:pos+1]...)
				res = append(res, log2[mergeado_hasta[1]:i+1]...)
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos

				mergeado_actual = true
			}

			pos = primeraAparicion(log1[mergeado_hasta[0]:], cambio.NombreCiudad)
			if pos != -1 {
				if !mergeado_actual {
					res = append(res, log2[mergeado_hasta[1]:i+1]...)
				}
				res = append(res, log1[mergeado_hasta[0]:pos+1]...)
				mergeado_hasta[0] = i
				mergeado_hasta[1] = pos
			}

			break
		case pb.InformanteReq_NUMBER_UPDATE:
			break
		}
	}

	// mergear cualquier cambio restante
	res = append(res, log1[mergeado_hasta[0]:]...)
	res = append(res, log2[mergeado_hasta[1]:]...)
	return res
}

func primeraAparicion(log []*pb.InformanteReq, ciudad string) int {
	for i, cambio := range log {
		if cambio.NombreCiudad == ciudad {
			return i
		}
	}
	return -1
}

func ultimaAparicion(log []*pb.InformanteReq, ciudad string) int {
	for i := len(log) - 1; i >= 0; i-- {
		cambio := log[i]
		if cambio.NombreCiudad == ciudad {
			return i
		}

		// en primeraAparicion no verificamos este caso porque implica algo que no debiera poder
		// ocurrir o que puede ser solucionado por una llamada sucesiva a mergeEntre2Logs
		if cambio.Comando == pb.InformanteReq_NAME_UPDATE && cambio.GetNuevaCiudad() == ciudad {
			return i
		}
	}
	return -1
}

func (server *FulcrumServer) obtenerLog() (*pb.MergeBeginRes, error) {
	res := &pb.MergeBeginRes{Reloj: &server.reloj}
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
		if err != nil {
			break
		}
		switch comando {
		case "AddCity":
			var rebeldes uint64
			fmt.Fscanf(registro_log, "%v\n", &rebeldes)
			res.Changelog = append(res.Changelog, &pb.InformanteReq{
				Comando:       pb.InformanteReq_ADD,
				NombrePlaneta: planeta,
				NombreCiudad:  ciudad,
				NuevoValor:    &pb.InformanteReq_NuevosRebeldes{NuevosRebeldes: rebeldes},
			})
			break
		case "DeleteCity":
			fmt.Fscanf(registro_log, "\n")
			res.Changelog = append(res.Changelog, &pb.InformanteReq{
				Comando:       pb.InformanteReq_DELETE,
				NombrePlaneta: planeta,
				NombreCiudad:  ciudad,
			})
			break
		case "UpdateValue":
			var rebeldes uint64
			fmt.Fscanf(registro_log, "%v\n", &rebeldes)
			res.Changelog = append(res.Changelog, &pb.InformanteReq{
				Comando:       pb.InformanteReq_NUMBER_UPDATE,
				NombrePlaneta: planeta,
				NombreCiudad:  ciudad,
				NuevoValor:    &pb.InformanteReq_NuevosRebeldes{NuevosRebeldes: rebeldes},
			})
			break
		case "UpdateName":
			var nueva_ciudad string
			fmt.Fscanf(registro_log, "%v\n", &nueva_ciudad)
			res.Changelog = append(res.Changelog, &pb.InformanteReq{
				Comando:       pb.InformanteReq_NAME_UPDATE,
				NombrePlaneta: planeta,
				NombreCiudad:  ciudad,
				NuevoValor:    &pb.InformanteReq_NuevaCiudad{NuevaCiudad: nueva_ciudad},
			})
			break
		default:
			fmt.Printf("comando desconocido: %v\n", comando)
			break
		}
	}

	// remover el log viejo
	os.Remove(archivo_log)

	return res, nil
}

func (server *FulcrumServer) integrarCambios(req *pb.MergeEndReq) error {
	for _, comando := range req.Changelog {
		var err error
		switch comando.Comando {
		case pb.InformanteReq_ADD:
			err = server.agregarCiudad(comando)
			break
		case pb.InformanteReq_NAME_UPDATE:
			err = server.cambiarNombre(comando)
			break
		case pb.InformanteReq_NUMBER_UPDATE:
			err = server.cambiarValor(comando)
			break
		case pb.InformanteReq_DELETE:
			err = server.borrarCiudad(comando)
			break
		}

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Actualizar los relojes a para reflejar los cambios
	for planeta := range server.planetas {
		server.planetas[planeta].X = req.Reloj.X
		server.planetas[planeta].Y = req.Reloj.Y
		server.planetas[planeta].Z = req.Reloj.Z
	}
	return nil
}

func (server *FulcrumServer) MergeBegin(ctx context.Context, req *pb.MergeBeginReq) (*pb.MergeBeginRes, error) {
	return server.obtenerLog()
}

func (server *FulcrumServer) MergeEnd(ctx context.Context, req *pb.MergeEndReq) (*pb.MergeEndRes, error) {
	err := server.integrarCambios(req)
	if err != nil {
		return nil, err
	}
	return &pb.MergeEndRes{}, nil
}

func max3(a int32, b int32, c int32) int32 {
	return max(a, max(b, c))
}

// NOTE(lucas): go no tiene un operador ternario asique esto tiene que ser hecho de la manera fea
func max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
