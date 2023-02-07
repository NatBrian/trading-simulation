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

# Run Program
`Make run`

# Challenge 1
## APIs:
1. API GET summary to get summary transactions per stock code
   `curl --location --request GET 'localhost:8080/summary'`
2. API POST upload to upload ndjson files inside a zip file `curl --location --request POST 'localhost:8080/upload' \  
   --form 'file=@"rawdata.zip"'`


## High level flow:
1. Admin upload all transaction for a day inside the rawdata.zip, by uploading via API POST upload
   - API returns `201` if there is no error in reading the zip file
   - Publish the topic to Kafka
2. Calculation OHLC is done asynchronously by the consumer
   - Store calculated OHLC to Redis by per stock codes
3. Admin gets all summaries via API GET summary
   - Get the summary in redis

Assumptions:
1. Type A is not used in Highest Price, Lowest Price, Open Price, and Closed Price
2. Transactions are uploaded only once per day, meaning that Redis is only storing transactions of that day.
   - Future improvements:
      - Redis TTL
      - Database to store a summary of each day, since DB is persistent and durable compared to Redis
      - In real life, transactions upload must be automated, hence the POST API upload need to be adjusted as needed
4. Admin will retrieve the latest summary via API GET summary, after the Kafka queue completed calculating all transactions.
   - Future improvements:
      - Async process to calculate OHLC can be improved by sending a notification to the user, when there is no uncommitted message / no lags in the kafka queue

# Challenge 2
Improvements:
1. Remove redundant stockCodes from the input parameter
2. Give proper variable names

# Feedbacks
1. Unclear problem statement.

*Every Transaction with type of A , E , and P will cause a change in the Stock’s OHLC with Volume & Value. Transaction that has Quantity = 0 is the Previous Price of a Stock. Previous Price is not accountable for Open Price, Highest Price, and Lowest Price of a Stock.*

*Every Transaction with type of E and P are accountable for Volume, Value, and Average Price of a Stock.*

From the above statements, I concluded that:
1. Type A, E, P are used for Prices (open, highest, lowest, close, and prev)
2. Type E P are used for Volume, Value, and Average Price

However, in the example, Highest Price = 8100
The example got 8100 because it does not calculate Type A with Price 8150.
After this misunderstanding, I asked the HR about this issue.
Her answer was: *A equals AddOrder which in ecommerce terms means ‘someone checked out an item but the seller has not confirmed the buy’, so it is not fair to include A to the trading volume.*
If so, then why is Open Price type A is correct? So, I reached my conclusion which I wrote above in *Challenge 1 assumption*.

2. Program requirements are unclear. I understand that the purpose of this exercise is to assess problem-solving skills. However, it is unclear if the desired outcome is a simple Golang Script, an App Server, or something else entirely. The desired complexity can vary based on candidates approach and implementation.