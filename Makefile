gen:
	mkdir -p pb
	protoc --go_out=pb --go-grpc_out=pb --proto_path=proto proto/*.proto 

clean: 
	rm ./pb/*.go
