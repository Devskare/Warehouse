package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name           string `env:"APP_NAME"`
	Production     bool   `env:"PRODUCTION"`
	GrpcServerPort string `env:"GRPC_SERVER_PORT"`
	Db             DB
}

type DB struct {
	Net      string `env:"DB_NET"`
	Driver   string `env:"DB_DRIVER"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	MaxConn  int    `env:"DB_MAXCONN"`
	Port     string `env:"DB_PORT"`
	Timeout  int    `env:"DB_TIMEOUT"`
}

func MustLoadConfig(env ...string) *AppConfig {
	var err error
	conf := &AppConfig{}
	err = godotenv.Load(env...)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(conf)
	if err != nil {
		panic(err)
	}

	return conf
}
