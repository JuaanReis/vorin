.PHONY: build clean rebuild run

BIN_NAME=vorin
OUT_DIR=./cmd
BIN_DIR=bin

build:
	go build -o $(BIN_NAME) 

clean:
	rm -rf $(BIN_NAME)

rebuild: clean build

run:
	./$(BIN_NAME) 