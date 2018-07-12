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

// Env defines data to be loaded from environment variables
var Env = struct {
	Port     int    `envconfig:"PORT" required:"true"`
	AppEnv   string `envconfig:"APP_ENV" default:"development"`
	MongoDSN string `envconfig:"MONGO_DSN" required:"true"`
}{}

func init() {
	err := envconfig.Process("RUBIXCORE", &Env)
	if err != nil {
		panic(err)
	}
}

func initLogger(target string) (*zap.Logger, error) {
	if target == "production" {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

func main() {
	logger, err := initLogger(Env.AppEnv)
	if err != nil {
		panic(err)
	}

	if Env.AppEnv == "development" {
		logger.Info("application configuration loaded successfully", zap.Any("appConfig", Env))
	}

	session, err := mgo.Dial(Env.MongoDSN)
	if err != nil {
		logger.Error("failed dialing mongo db connection", zap.Error(err))
		panic(err)
	}

	err = db.InitDB(session)
	if err != nil {
		logger.Error("failed initializing db", zap.Error(err))
		panic(err)
	}

	routes := api.InitRoutes(session, logger)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", Env.Port),
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           routes,
	}

	logger.Info("Server listening on : ", zap.Any("url", fmt.Sprintf("http://0.0.0.0:%d", Env.Port)))
	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}
