gen-server:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	calculator/calculatorpb/calculator.proto

run-server: 
	go run calculator/server/server.go

run-client:
	go run calculator/client/client.go
