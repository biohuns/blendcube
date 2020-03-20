APP_NAME := blendcube

default: build

.PHONY: build
build:
	go build -o $(APP_NAME)

.PHONY: run
run: build
	./$(APP_NAME)

.PHONY: clean
clean:
	rm -f $(APP_NAME)
