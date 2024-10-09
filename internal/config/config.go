package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Add string `yaml:"add" env:"HTTP_SERVER_ADD" env-required:"true"`
}

// Config struct represents the configuration for the application
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

// MustLoad reads the configuration from a file or environment variables
func MustLoad() *Config {
	var configpath string

	// Check for CONFIG_PATH environment variable
	configpath = os.Getenv("CONFIG_PATH")
	if configpath == "" {
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()

		configpath = *flags

		if configpath == "" {
			log.Fatal("config path not set")
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configpath)
	}

	var cfg Config

	// Read the configuration from the file
	err := cleanenv.ReadConfig(configpath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cfg // Return the loaded configuration
}
