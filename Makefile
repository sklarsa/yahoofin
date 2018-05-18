build:
	go build .
	go build -o bin/yahoofin github.com/sklarsa/yahoofin/yahoofin-cli

test:
	go test

all: build test
