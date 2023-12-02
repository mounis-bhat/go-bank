.PHONY: build run test

BUILD_DIR := ./bin
APP_NAME := gobank
MAIN_FILE := cmd/app/main.go

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	@$(BUILD_DIR)/$(APP_NAME)

test:
	@go test -v ./...

clean:
	@rm -rf $(BUILD_DIR)
