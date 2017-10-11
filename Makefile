build:
	go build ./

test:
	go test ./

all: build test
	go build ${LDFLAGS} -o ${BINARY} ./
