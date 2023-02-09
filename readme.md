
# Trading Simulation

  

Requirements:

1. Receive lots of transactions per day

2. Calculate Open-high-low-close (OHLC) summary for each stock code

  

# High Level Flow

- Transactions per day send to App Server

- App server produce transactions to Kafka Queue

- Consumer consume those transactions to calculate OHLC asynchronously, then store it in Redis

- Admin can get summary via API

  

# Implementations

## APIs

1. GET SUMMARY: get summary of calculated OHLC for all stock codes

```
curl --location --request GET 'localhost:8080/summary'
```

2. POST UPLOAD: upload transaction ndjson files inside a zip file

```
curl --location --request POST 'localhost:8080/upload' \

--form 'file=@"rawdata.zip"'
```

*example file is rawdata.zip*

## Kafka

- Publish transactions

- Consume transactions to calculate OHLC

  

## Redis

- Store calculation of OHLC based on stock code

  

## Protobuf

  

# Future Improvements

1. Feeding transactions data to App Server

2. Notification when calculation is finished, so Admin does not need to manually hit GET API to app server

  

# Installation

  

## Kafka

1. `cd kafka`

2. `docker-compose up -d`

3. `docker exec broker \

kafka-topics --bootstrap-server broker:9092 \

--create \

--topic process_change_record`

  

## Redis

1. `docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest`

  

## Protobuf

1. `brew install protobuf`

2. `go install google.golang.org/grpc`

3. `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

4. `protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative proto/*.proto`

  

# Run Program
`Make run`