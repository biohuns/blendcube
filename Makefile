APP_NAME := blendcube

default: build

.PHONY: build
build:
ifeq ($(OS),Windows_NT)
	@go build -o $(APP_NAME).exe
else
	@go build -o $(APP_NAME)
endif

.PHONY: run
run: build
ifeq ($(OS),Windows_NT)
	@./$(APP_NAME).exe
else
	@./$(APP_NAME)
endif

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: clean
clean:
	go clean -x
