package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// App defines shared dependencies, request handlers,
// and url-mappings for the API
type App struct {
	session  *mgo.Session
	logger   *zap.Logger
	upgrader *websocket.Upgrader
}

// Router returns a http.Handler with url mappings
// for all routes handlers in the API
func (a *App) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/queues", a.listQueues)
	r.Post("/queues", a.createQueue)

	r.Get("/ws/status", a.handleStatusCheck)
	r.Get("/ws/test", a.handleStatusTest)

	fileServer(r, "/static", http.Dir("./ui/static"))

	return r
}

// NewApp returns a pointer to a new app with session and logger
// and websocket upgrader properly configured and ready for use in all routes
func NewApp(s *mgo.Session, l *zap.Logger, u *websocket.Upgrader) *App {
	a := App{s, l, u}

	return &a
}
