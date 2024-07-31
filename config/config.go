package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		MySQL  MySQL
		Server Server
		IsDev  bool `envconfig:"IS_DEV"`
	}

	MySQL struct {
		Name     string `envconfig:"MYSQL_NAME"`
		Host     string `envconfig:"MYSQL_HOST"`
		User     string `envconfig:"MYSQL_USER"`
		Password string `envconfig:"MYSQL_PASSWORD"`
	}

	Server struct {
		Port string `envconfig:"SERVER_PORT"`
	}
)

var Cfg Config

func FillConfig() error {
	return envconfig.Process("", &Cfg)
}
