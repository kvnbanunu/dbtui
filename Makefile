TARGET = bin/dbtui
MAIN = main.go
DB_PATH = ./sqlite.db
FLAGS = $(DB_PATH)

build: clean
	go build $(MAIN) -o $(TARGET)

run:
	go run $(MAIN) $(FLAGS)

seed:
	go run $(MAIN) -seed $(FLAGS)

clean:
	rm -rf bin
