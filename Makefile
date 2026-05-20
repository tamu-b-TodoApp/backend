ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: build run dev clean

build:
	go build -o server .

run:
	./server

dev:
	go run main.go

clean:
	rm -f server

migrate-apply:
	atlas migrate apply --env $(ENV)

migrate-diff:
	atlas migrate diff $(name) --env $(ENV)