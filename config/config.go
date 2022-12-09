package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Songmu/replaceablewriter"
	"github.com/go-chi/chi/middleware"
)

type Config struct {
	Server struct {
		Port    int           `json:"port"`
		Timeout time.Duration `json:"timeout"`
	} `json:"server"`
	Log struct {
		FilePath string `json:"file_path"`
	} `json:"log"`
	Model struct {
		FilePath       string `json:"file_path"`
		BinaryFilePath string `json:"binary_file_path"`
	} `json:"model"`
}

var (
	Shared  = new(Config)
	logFile *replaceablewriter.Writer
)

func (c *Config) GetPort() string {
	return ":" + strconv.Itoa(c.Server.Port)
}

func Configure(exit chan int) error {
	configPath := flag.String(
		"config",
		"./config.json",
		"config filepath",
	)
	flag.Parse()

	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, Shared); err != nil {
		return err
	}
	if err := configureLogger(exit); err != nil {
		return err
	}
	configureReopenLogFile(exit)

	return nil
}

func configureLogger(exit chan int) error {
	if Shared.Log.FilePath != "" {
		var err error
		logFile, err = openLogFile()
		if err != nil {
			return err
		}
		middleware.DefaultLogger = middleware.RequestLogger(
			&middleware.DefaultLogFormatter{
				Logger:  log.New(logFile, "", log.LstdFlags),
				NoColor: true,
			},
		)

		configureLogger(exit)
	}

	return nil
}

func configureReopenLogFile(exit chan int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGHUP:
				if err := reopenLogFile(); err != nil {
					log.Printf("reopen log error: %s", err)
					exit <- 1
				}
				log.Println("reopen log")
			case syscall.SIGTERM:
				log.Println("shutdown...")
				exit <- 0
				return
			}
		}
	}()
}

func reopenLogFile() error {
	if err := logFile.Close(); err != nil {
		return err
	}
	f, err := openLogFile()
	if err != nil {
		return fmt.Errorf("failed to open logfile: %s", err)
	}
	logFile.Replace(f)
	return nil
}

func openLogFile() (*replaceablewriter.Writer, error) {
	f, err := os.OpenFile(
		Shared.Log.FilePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return replaceablewriter.New(f), nil
}
