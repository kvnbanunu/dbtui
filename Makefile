TARGET = bin/main
MAIN = cmd/main.go

all: clean build

build:
	go build $(MAIN)

run:
	go run $(MAIN)

clean:
	rm -rf bin/
