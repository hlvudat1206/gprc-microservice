gen-cal:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	calculator/calculatorpb/calculator.proto

gen-contact:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	contact/contactpb/contact.proto

run-server: 
	go run calculator/server/server.go

run-client:
	go run calculator/client/client.go

run-contact-server: 
	go run contact/server/server.go contact/server/models.go

run-contact-client:
	go run contact/client/client.go
