package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rubixq/rubixcore/pkg/api"
	"github.com/rubixq/rubixcore/pkg/db"
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

	session, err := mgo.Dial(env.MongoDSN)
	if err != nil {
		logger.Error("failed dialing mongo db connection", zap.Any("error", err))
		panic(err)
	}

	err = db.InitDB(session)
	if err != nil {
		logger.Error("failed initializing db", zap.Any("error", err))
		panic(err)
	}

	routes := api.InitRoutes(session, logger)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", env.Port),
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           routes,
	}

	logger.Info("Server listening on : ", zap.Int("port", env.Port))
	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}
