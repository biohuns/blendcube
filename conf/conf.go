package conf

import (
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
			Port    int           `toml:"port"`
			Timeout time.Duration `toml:"timeout"`
		} `toml:"server"`
		Log struct {
			StdOutput  bool   `toml:"std_output"`
			FileOutput bool   `toml:"file_output"`
			FilePath   string `toml:"file_path"`
			NoColor    bool   `toml:"no_color"`
		} `toml:"log"`
		Model struct {
			FilePath string `toml:"file_path"`
		}
	}
)

var (
	Shared *Config
)

// GetPort
func (c *Config) GetPort() string {
	return ":" + strconv.Itoa(c.Server.Port)
}

// Configure
func Configure() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.SetConfigFile("./config.toml")
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
