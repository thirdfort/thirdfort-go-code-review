package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Environment string

const (
	Local       Environment = "local"
	Development Environment = "development"
	Testing     Environment = "testing"
	Sandbox     Environment = "sandbox"
	Production  Environment = "production"

	configFile string = "./config.yaml"
)

type Config struct {
	App struct {
		Name       string      `env-default:"consumer-api"`
		Env        Environment `env:"ENV" env-default:"local"`
		Debug      bool        `env:"DEBUG" env-default:"true"`
		NoAuth     bool        `env:"NO_AUTH" env-default:"false"` // Skip local authentication
		SofEnabled bool        `env:"SOF_ENABLED" env-default:"false"`
	} `yaml:"app"`
	Service struct {
		Port    int    `env:"PORT" env-default:"6060"`
		Address string `env:"ADDRESS" env-default:"0.0.0.0"`
	} `yaml:"service" env-prefix:"CONAPI_"`
	Postgres struct {
		User        string `env:"USER" yaml:"user" env-default:"postgres"`
		Password    string `env:"PASSWORD" yaml:"password" env-default:"postgres"`
		Host        string `env:"HOST" yaml:"host" env-default:"postgres"`
		Database    string `env:"NAME" yaml:"name" env-default:"consumer-api"`
		Port        string `env:"PORT" yaml:"port" env-default:"5432"`
		EnableDebug bool   `env:"DEBUG" yaml:"debug" env-default:"false"`
	} `yaml:"postgres" env-prefix:"CONAPI_POSTGRES_"`
	Clients struct {
		PlatformAPI struct {
			Aud       string `env:"AUD"`
			Url       string `env:"URL"`
			TimeoutMs int    `env:"TIMEOUT" yaml:"timeout_ms" env-default:"15000"`
			Retries   int    `env:"RETRYMAX" yaml:"retry_max" env-default:"5"`
		} `yaml:"platform-api" env-prefix:"PA_"`
	} `yaml:"clients" env-prefix:"CONAPI_"`
}

func New() (*Config, error) {
	var conf Config

	env := os.Getenv("ENV")

	_, err := os.Stat(configFile)
	if errors.Is(err, os.ErrNotExist) || env == string(Testing) {
		err = cleanenv.ReadEnv(&conf)
		if err != nil {
			return nil, errors.Wrap(err, "ReadEnv: error processing envconfig")
		}

		// if conf.App.Debug {
		// 	printConf(conf)
		// }

		return &conf, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "Stat: error processing envconfig")
	}

	err = cleanenv.ReadConfig(configFile, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig: error processing envconfig")
	}

	if conf.App.Debug {
		// printConf(conf)
	}

	return &conf, nil
}
