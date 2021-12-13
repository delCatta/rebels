gen:
	mkdir -p pb
	protoc --go_out=pb --go-grpc_out=pb --proto_path=proto proto/*.proto 
i:
	go run ./cmd/informante/main.go	
l:
	go run ./cmd/leia/main.go
b:
	go run ./cmd/broker/main.go	
f2:
	go run ./cmd/fulcrum/main.go X
f3:
	go run ./cmd/fulcrum/main.go Z
f4:
	go run ./cmd/fulcrum/main.go Y
	
clean: 
	rm ./pb/*.go
