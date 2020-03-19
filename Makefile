APP_NAME := blendcube

default: build

.PHONY: build
build:
	go build -o BUILD/$(APP_NAME)

.PHONY: run
run: build
	./BUILD/$(APP_NAME)

.PHONY: clean
clean:
	rm -f BUILD/$(APP_NAME)
