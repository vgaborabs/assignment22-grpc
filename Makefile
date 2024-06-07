start:
	@go build -o bin/usergrpc ./cmd/main.go
	@./bin/usergrpc
build:
	@docker build -t assignment22-grpc .
run:
	@docker run -d -p 5000:5000 assignment22-grpc:latest
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto