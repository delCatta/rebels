gen:
	mkdir -p pb
	protoc --go_out=pb --go-grpc_out=pb --proto_path=proto proto/*.proto 
i:
	go run ./cmd/informante/main.go	

clean: 
	rm ./pb/*.go
