APP_NAME := blendcube

default: build

.PHONY: build
build:
	@GO111MODULE=on go build -o $(APP_NAME)

.PHONY: build-win
build-win:
	@set GO111MODULE=on
	@go build -o $(APP_NAME)

.PHONY: run
run: build
	@./$(APP_NAME)

.PHONY: run-win
run-win: build-win
	@./$(APP_NAME)

.PHONY: clean
clean:
	GO111MODULE=on go clean -x
