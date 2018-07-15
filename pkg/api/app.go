package api

import (
	"net/http"
	"sync"

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
	lock     *sync.Mutex
	kiosks   map[string]*websocket.Conn
	counters map[string]*websocket.Conn
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
	a := App{
		session:  s,
		logger:   l,
		upgrader: u,
	}

	a.kiosks = make(map[string]*websocket.Conn)
	a.counters = make(map[string]*websocket.Conn)

	return &a
}

func (a *App) addKiosk(id string, conn *websocket.Conn) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.kiosks[id] = conn
}

func (a *App) removeKiosk(id string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.kiosks, id)
}

func (a *App) addCounter(id string, conn *websocket.Conn) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.counters[id] = conn
}

func (a *App) removeCounter(id string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.counters, id)
}
