package main

import (
	"blendcube/conf"
	"blendcube/cube"
	"blendcube/handler"
	"log"
	"net/http"
)

func main() {
	if err := start(); err != nil {
		log.Fatalln(err)
	}
}

func start() error {
	if err := conf.Configure(); err != nil {
		return err
	}

	if err := cube.Initialize(); err != nil {
		return err
	}

	if err := http.ListenAndServe(
		conf.Shared.GetPort(),
		handler.New(),
	); err != nil {
		return err
	}

	return nil
}
