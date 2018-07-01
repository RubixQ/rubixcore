package main

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type config struct {
	Port   int    `envconfig:"PORT" required:"true"`
	AppEnv string `envconfig:"APP_ENV" default:"development"`
	SQLDSN string `envconfig:"SQL_DSN" required:"true"`
}

func loadConfig() (*config, error) {
	c := new(config)
	err := envconfig.Process("rubixcore", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	env, err := loadConfig()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger
	if env.AppEnv == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	logger.Info("Application configuration loaded successfully", zap.Any("appConfig", env))
}
