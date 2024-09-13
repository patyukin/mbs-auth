.PHONY: start stop rebuild gen up down restart deps tidy

up:
	docker compose up -d

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop

restart:
	docker compose restart

rebuild:
	docker compose down -v --remove-orphans
	docker compose up -d --build

gen:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api_gateway/main.go -o docs/

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

tidy:
	go mod tidy
