package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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
	} `yaml:"service" env-prefix:"API_"`
	Postgres struct {
		User        string `env:"USER" yaml:"user" env-default:"postgres"`
		Password    string `env:"PASSWORD" yaml:"password" env-default:"postgres"`
		Host        string `env:"HOST" yaml:"host" env-default:"postgres"`
		Database    string `env:"NAME" yaml:"name" env-default:"consumer-api"`
		Port        string `env:"PORT" yaml:"port" env-default:"5432"`
		EnableDebug bool   `env:"DEBUG" yaml:"debug" env-default:"false"`
	} `yaml:"postgres" env-prefix:"API_POSTGRES_"`
	Clients struct {
		Onfido struct {
			AppID struct {
				Android string `env:"ANDROID" yaml:"android" env-default:"app.thirdfort.android.sandbox"`
				Ios     string `env:"IOS" yaml:"ios" env-default:"app.thirdfort.ios.sandbox"`
			} `env:"APPID_" yaml:"app_id"`
		} `yaml:"onfido" env-prefix:"API_ONFIDO_"`
		PlatformAPI struct {
			Aud       string `env:"AUD" env-default:"https://api-t.thirdfort.io/v2/"`
			Url       string `env:"URL" env-default:"https://sandbox.api.thirdfort.io"`
			TimeoutMs int    `env:"TIMEOUT" yaml:"timeout_ms" env-default:"15000"`
			Retries   int    `env:"RETRYMAX" yaml:"retry_max" env-default:"5"`
		} `yaml:"platform-api" env-prefix:"PA_"`
		Jwt struct {
			Aud     []string `env:"AUD" env-default:"[https://dev.api.thirdfort.io/v2,sandbox-thirdfort-app,https://sandbox.api.thirdfort.io/v2/,https://sandbox-tf.eu.auth0.com/userinfo]"`
			Issuers struct {
				Consumer []string `env:"ISSUERS_CONSUMER" env-default:"[https://securetoken.google.com/sandbox-thirdfort-app]"`
			}
		} `yaml:"jwt" env-prefix:"JWT_"`
	} `yaml:"clients" env-prefix:"API_"`
	Otel struct {
		Enabled         bool          `env:"ENABLED" env-default:"false"`
		PushInterval    time.Duration `env:"PUSH_INTERVAL" yaml:"push_interval" env-default:"10000ms"`
		RuntimeInterval time.Duration `env:"RUNTIME_INTERVAL" yaml:"runtime_interval" env-default:"15000ms"`
		SamplingRatio   float64       `env:"SAMPLING_RATIO" yaml:"sampling_ratio" env-default:"1"`
	} `yaml:"otel" env-prefix:"API_OTEL_"`
	Sentry struct {
		Dsn         string      `env:"DSN" yaml:"dsn" env-default:""`
		Environment Environment `env:"ENV" yaml:"environment" env-default:"local"`
	} `yaml:"sentry" env-prefix:"API_SENTRY_"`
	PubSub struct {
		ProjectID string `env:"PROJECT_ID" env-default:"local"`
		Timeout   int    `env:"TIMEOUT" env-default:"500"`
		Topic     string `env:"TOPIC" env-default:"consumer-api"`
	} `yaml:"pubsub" env-prefix:"API_PUBSUB_"`
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

func printConf(conf Config) {
	c, err := json.MarshalIndent(conf, "", "  ")
	if err == nil {
		fmt.Printf("Config: %s\n", string(c))
	}
}
