package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
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

	upgrader := &websocket.Upgrader{}
	app := api.NewApp(session, logger, upgrader)
	router := app.Router()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", Env.Port),
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router),
	}

	// Run server in a goroutine so that it doesn't block.
	go func() {
		logger.Info("api accessible from : ", zap.Any("url", fmt.Sprintf("http://0.0.0.0:%d", Env.Port)))
		logger.Info("ws accessible from : ", zap.Any("url", fmt.Sprintf("http://0.0.0.0:%d/ws", Env.Port)))
		logger.Info("ws test accessible from : ", zap.Any("url", fmt.Sprintf("http://0.0.0.0:%d/ws/test", Env.Port)))
		logger.Info("ws status accessible from : ", zap.Any("url", fmt.Sprintf("http://0.0.0.0:%d/ws/status", Env.Port)))

		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	// Perform graceful shutdowns when quit via SIGINT (Ctrl+C)
	// or SIGKILL, SIGQUIT or SIGTERM (Ctrl+/)
	signal.Notify(ch, os.Interrupt, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)

	// Block until signal is received.
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	logger.Info("shutting down server")
	os.Exit(0)

}
