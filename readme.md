# Installation

### golangci-lint
1. go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.0

### kafka
1. `cd kafka`
2. `docker-compose up -d`
3. `docker exec broker \
   kafka-topics --bootstrap-server broker:9092 \
   --create \
   --topic process_change_record`

### redis
1. `docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest`

### protobuff
1. `brew install protobuf`
2. `go install google.golang.org/grpc`
3. `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
4. `protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative proto/*.proto`