package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/rubixq/rubixcore/pkg/api/controllers"
	"go.uber.org/zap"
	redis "gopkg.in/redis.v4"
)

// Resources groups dependencies needed by handler funcs
type Resources struct {
	db       *sqlx.DB
	logger   *zap.Logger
	redis    *redis.Client
	upgrader *websocket.Upgrader
}

// NewResources returns a pointer to a group of dependencies for handlers
func NewResources(db *sqlx.DB, logger *zap.Logger, redis *redis.Client, upgrader *websocket.Upgrader) *Resources {
	return &Resources{
		db:       db,
		logger:   logger,
		redis:    redis,
		upgrader: upgrader,
	}
}

// Router returns a http.Handler configured with route handlers
func Router(res *Resources) http.Handler {
	r := chi.NewRouter()

	r.Post("/auth", controllers.Authenticate(res.db, res.logger))
	return handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)
}
