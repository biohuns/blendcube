package main

import (
	"blendcube/conf"
	"blendcube/cube"
	"blendcube/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := conf.Configure(); err != nil {
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
		Addr:    conf.Shared.GetPort(),
		Handler: handler.New(),
	}

	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan int)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			exitChan <- 1
		}
	}()

	go func() {
		s := <-signalChan
		switch s {
		case syscall.SIGHUP, syscall.SIGTERM:
			log.Println("shutdown...")
			exitChan <- 0
		default:
			log.Printf("receive unknown signal: %+v\n", s)
			exitChan <- 1
		}
	}()

	exitCode := <-exitChan
	os.Exit(exitCode)
}
