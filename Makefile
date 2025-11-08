TARGET = bin/main
MAIN = main.go

all: clean build

build:
	go build $(MAIN) -o $(TARGET)

run:
	go run $(MAIN)

clean:
	rm -rf bin
