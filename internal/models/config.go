package models

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// Config struct represents the configuration for the application
// env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}
