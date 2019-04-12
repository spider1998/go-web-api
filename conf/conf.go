package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug    bool   `json:"debug" default:"false"`
	HTTPAddr string `json:"http_addr" default:":80"`
	Mysql    string `json:"mysql" default:""`
	Redis    string `json:"redis" default:"localhost:6379"`
	Version  string `json:"version" default:"0.0.1"`
}

func NewConfig() (Config, error) {
	godotenv.Load()
	var config Config
	err := envconfig.Process("", &config)
	return config, err
}
