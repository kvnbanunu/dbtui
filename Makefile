TARGET = bin/main
MAIN = main.go
DB_PATH = ./sqlite.db
FLAGS = -seed $(DB_PATH)

all: clean build

build:
	go build $(MAIN) -o $(TARGET)

run:
	go run $(MAIN) $(FLAGS)

clean:
	rm -rf bin
