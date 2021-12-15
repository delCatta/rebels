# Sistemas Distribuídos Grupo 39

## Integrantes:
- Felipe Cisternas 201873022-K
- Diego Cattarnich 201873023-8
- Lucas Galindo 201873004-1


# Instrucciones de Uso.

En todas las máquinas lo primero que se debe hacer es correr `make gen` para generar los archivos de grpc.



## Máquina 1: Broker

Instrucciones:

- Correr `make gen`.
- Iniciar el Broker `make b`


## Máquina 2: Leia y Fulcrum
Instrucciones:
- Correr `make gen`.
- Inicializar Fulcrum `make f2`.
- Para ejecutar Leia: `make l`.


## Máquina 3: Fulcrum
Instrucciones:
- Correr `make gen`.
- Inicializar Fulcrum `make f3`.

  
## Máquina 4: Informantes y Fulcrum
Instrucciones:
- Correr `make gen`.
- Inicializar Fulcrum `make f4`.
- Para ejecutar Informante `make i`.


*Nota Importante:* lo expuesto anteriormente son las instrucciones por máquina,
el orden en que se deben instanciar los procesos es:
```
make f2 y make f4 -> make f3 -> make b -> make i y make l
```

## Consideraciones sobre el Merge

La manera en que resuelve conflictos el merge es que determina un orden lógico en el que deben ocurrir los cambios sobre un archivo ocupando como referencia las instrucciones AddCity, UpdateName y DeleteCity (ej: si encuentra un AddCity en el log entonces se determina que la primera mención de esa ciudad en los otros ocurre inmediatamente despues. Realiza un proceso similar con UpdateName y DeleteCity), el problema ocurre si el conflicto es entre cual UpdateNumber ocurre primero ya que no hay una buena manera de determinarlo. Lo que hacemos para resolver este caso es que simplemente se coloca primero el que pertenezca al log del primer parámetro de la función (la función que hace el merge hace llamadas recursivas a si misma donde se cambia el orden de los logs por lo que no necesariamente es el log del primer fulcrum).
