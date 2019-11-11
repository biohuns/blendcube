package conf

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

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
			StdOutput  bool
			FileOutput bool
			FilePath   string
			NoColor    bool
		}
		Model struct {
			FilePath string
		}
	}
)

var (
	// Shared holds config instance
	Shared *Config
	conf   = flag.String("conf", "./config.toml", "config file path")
)

// GetPort return port string
func (c *Config) GetPort() string {
	return ":" + strconv.Itoa(c.Server.Port)
}

// Configure parse config file and environment variable
func Configure() error {
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

	return nil
}

func configureLogger() error {
	var writer []io.Writer

	if Shared.Log.StdOutput {
		writer = append(writer, os.Stdout)
	}
	if Shared.Log.FileOutput {
		Shared.Log.NoColor = true
		logFile, err := os.OpenFile(
			Shared.Log.FilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return err
		}
		writer = append(writer, logFile)
	}

	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger: log.New(
				io.MultiWriter(writer...),
				"",
				log.LstdFlags),
			NoColor: Shared.Log.NoColor,
		},
	)

	return nil
}
