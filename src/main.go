package main

import (
	"blendcube/config"
	"blendcube/cube"
	"blendcube/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	exit := make(chan int)

	if err := config.Configure(exit); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("server configure: success")

	if err := cube.Initialize(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("loading model: success")

	srv := &http.Server{
		Addr:    config.Shared.GetPort(),
		Handler: handler.New(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			exit <- 1
		}
	}()

	exitCode := <-exit
	os.Exit(exitCode)
}
