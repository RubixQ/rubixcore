package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/rubixq/rubixcore/api"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

type config struct {
	Port     int    `envconfig:"PORT" required:"true"`
	AppEnv   string `envconfig:"APP_ENV" default:"development"`
	MongoDSN string `envconfig:"MONGO_DSN" required:"true"`
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

	logger.Info("application configuration loaded successfully", zap.Any("appConfig", env))

	_, err = mgo.Dial(env.MongoDSN)
	if err != nil {
		logger.Error("failed dialing mongo db connection", zap.Any("error", err))
		panic(err)
	}

	r := api.InitRoutes()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: r,
	}

	logger.Info("Server listening on : ", zap.Int("port", env.Port))
	log.Panic(s.ListenAndServe())

}
