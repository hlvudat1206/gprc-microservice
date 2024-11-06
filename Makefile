gen-cal:
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
