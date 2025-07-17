.PHONY: build clean rebuild install

BIN_NAME=vorin
OUT_DIR=./cmd
BIN_DIR=bin

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) $(OUT_DIR)

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)

rebuild: clean build

install:
	bash cmd/script/install.sh