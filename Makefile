setup:


run:
	set -o allexport; source config/config.env; set +o allexport && go run main.go

lint:
	golangci-lint run