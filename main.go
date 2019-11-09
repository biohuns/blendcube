package main

import (
	"blendcube/conf"
	"blendcube/cube"
	"blendcube/handler"
	"log"
	"net/http"
)

func main() {
	if err := conf.Configure(); err != nil {
		log.Fatalln(err)
	}

	if err := cube.Initial(); err != nil {
		log.Fatalln(err)
	}

	if err := http.ListenAndServe(
		conf.Shared.GetPort(),
		handler.New(),
	); err != nil {
		log.Fatalln(err)
	}
}
