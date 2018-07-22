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
	session    *mgo.Session
	logger     *zap.Logger
	upgrader   *websocket.Upgrader
	counters   map[string]*websocket.Conn
	nextTicket int
}

// Router returns a http.Handler with url mappings
// for all routes handlers in the API
func (a *App) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/queues", a.listQueues)
	r.Post("/queues", a.createQueue)

	r.Get("/kiosks/new", a.handleKioskSetup)

	r.Get("/ws/status", a.handleStatusCheck)
	r.Get("/ws/test", a.handleStatusTest)

	r.Post("/customers", a.createCustomer)

	fileServer(r, "/static", http.Dir("./ui/static"))

	return r
}

// NewApp returns a pointer to a new app with session and logger
// and websocket upgrader properly configured and ready for use in all routes
func NewApp(s *mgo.Session, l *zap.Logger, u *websocket.Upgrader) *App {
	a := App{
		session:    s,
		logger:     l,
		upgrader:   u,
		counters:   make(map[string]*websocket.Conn),
		nextTicket: 0,
	}

	return &a
}

func (a *App) addCounter(id string, conn *websocket.Conn) {
	a.counters[id] = conn
}

func (a *App) getCounterConn(id string) *websocket.Conn {
	return a.counters[id]
}
