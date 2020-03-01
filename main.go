package main

import (
	"log"
	"os"

	"github.com/biohuns/blendcube/config"
	"github.com/biohuns/blendcube/cube"
	"github.com/biohuns/blendcube/handler"
)

var exit = make(chan int)

func start() {
	if err := config.Configure(exit); err != nil {
		log.Fatalln(err)
	}
	log.Println("server configure: success")

	if err := cube.Initialize(); err != nil {
		log.Fatalln(err)
	}
	log.Println("loading model: success")

	go func() {
		if err := handler.NewServer().ListenAndServe(); err != nil {
			log.Println(err)
			exit <- 1
		}
	}()

	exitCode := <-exit
	os.Exit(exitCode)
}

func main() {
	start()
}
