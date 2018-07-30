package main

import (
	"context"
	"fmt"
	"net"
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
	"gopkg.in/redis.v4"
)

// Env defines data to be loaded from environment variables
var Env = struct {
	Port                int    `envconfig:"PORT" required:"true"`
	AppEnv              string `envconfig:"APP_ENV" default:"development"`
	MongoDSN            string `envconfig:"MONGO_DSN" required:"true"`
	RedisURL            string `envconfig:"REDIS_URL" required:"true"`
	TicketResetInterval int    `envconfig:"TICKET_RESET_INTERVAL" required:"true"`
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

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	logger.Info("redis connection setup successfully", zap.Any("pong", pong))

	upgrader := &websocket.Upgrader{}
	app := api.NewApp(session, client, logger, upgrader)
	router := app.Router()

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", Env.Port))
	if err != nil {
		logger.Fatal("failed binding to port", zap.Int("port", Env.Port))
	}
	defer listener.Close()

	url := fmt.Sprintf("http://%s", listener.Addr())
	logger.Info("server listening on ", zap.String("url", url))

	server := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router),
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		recv := <-sigs
		logger.Info("received signal, shutting down", zap.Any("signal", recv.String))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Warn("error shutting down server", zap.Error(err))
		}
	}()

	// Sets a timer to reset the ticket number generation after n hours
	// as specified in Env.TicketResetInterval
	ticker := time.NewTicker(time.Duration(Env.TicketResetInterval) * time.Hour)
	go func() {
		for range ticker.C {
			logger.Info("resetting ticket numbering")
			app.ResetTicketing()
		}
	}()

	if err := server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			logger.Fatal("http.Serve returned an error",
				zap.Error(err),
			)
		}
	}

	// Waits for all idle connections to be closed during shutdown
	<-idleConnsClosed
	logger.Info("server shutdown successfully")
}
