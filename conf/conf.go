package conf

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Songmu/replaceablewriter"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server struct {
			Port    int
			Timeout time.Duration
		}
		Log struct {
			Output   string
			FilePath string
		}
		Model struct {
			FilePath       string
			BinaryFilePath string
		}
	}
)

var (
	conf = flag.String("conf", "./config.toml", "config file path")

	// Shared holds config instance
	Shared *Config
	// LogFile holds log output file
	LogFile *replaceablewriter.Writer
)

// GetPort return port string
func (c *Config) GetPort() string {
	return ":" + strconv.Itoa(c.Server.Port)
}

// Configure parse config file and environment variable
func Configure(exit chan int) error {
	flag.Parse()

	viper.GetString("config_path")
	viper.SetConfigType("toml")
	viper.SetConfigFile(*conf)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Shared); err != nil {
		return err
	}

	if err := configureLogger(); err != nil {
		return err
	}
	configureReopenLogFile(exit)

	return nil
}

func configureLogger() error {
	switch Shared.Log.Output {
	case "stdout":
		middleware.DefaultLogger = middleware.RequestLogger(
			&middleware.DefaultLogFormatter{
				Logger: log.New(os.Stdout, "", log.LstdFlags),
			},
		)
	case "file":
		f, err := os.OpenFile(
			Shared.Log.FilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return err
		}
		LogFile = replaceablewriter.New(f)
		middleware.DefaultLogger = middleware.RequestLogger(
			&middleware.DefaultLogFormatter{
				Logger:  log.New(LogFile, "", log.LstdFlags),
				NoColor: true,
			},
		)
	default:
		return errors.New("output must be stdout or file")
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
		for {
			s := <-sig
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
			default:
				log.Printf("receive unknown signal: %+v\n", s)
				exit <- 1
				return
			}
		}
	}()
}

func reopenLogFile() (err error) {
	if err := LogFile.Close(); err != nil {
		return fmt.Errorf("logfile close error: %s", err)
	}
	f, err := os.OpenFile(
		Shared.Log.FilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("logfile open error: %s", err)
	}
	LogFile.Replace(f)
	return nil
}
