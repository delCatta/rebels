package service

import (
	"github.com/delCatta/rebels/pb"
)

type RegistroPlanetario struct {
	ciudades map[string][]*pb.InformanteReq
}

func merge(logs [3]*pb.MergeInfo) (*pb.MergeInfo, error) {
	merge_result := &pb.MergeInfo{}

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
