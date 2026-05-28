ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: build run dev clean seed gen-jwt-secret test

test:
	go test ./...

build:
	go build -o server .

run:
	./server

dev:
	go run main.go

clean:
	rm -f server

seed:
	go run ./cmd/seed/main.go

gen-jwt-secret:
	go run ./cmd/gen-jwt-secret/main.go

migrate-apply:
	atlas migrate apply --env $(ENV)

migrate-diff:
	atlas migrate diff $(name) --env $(ENV)

migrate-reset:
	atlas schema clean --url "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" --auto-approve
	atlas migrate apply --env $(ENV)