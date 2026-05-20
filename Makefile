.PHONY: build run dev clean

build:
	go build -o server .

run:
	./server

dev:
	go run main.go

clean:
	rm -f server
