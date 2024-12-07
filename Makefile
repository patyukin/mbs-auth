.PHONY: start stop rebuild gen up down restart deps tidy gen-auth-api vendor-proto

LOCAL_BIN:=$(CURDIR)/bin

tidy:
	GOBIN=$(LOCAL_BIN) go mod tidy

test:
	go test ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...