package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/redis.v4"
)

var (
	lastResetTimestamp              prometheus.Gauge
	failedCustomerRegistrations     prometheus.Counter
	successfulCustomerRegistrations prometheus.Counter
)

func init() {
	lastResetTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "last_reset_timestamp",
		Help: "The Unix timestamp of when last ticket generation algorithm was reset",
	})

	failedCustomerRegistrations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_customer_registration_count",
		Help: "The number of failed customer registration attempts since process start",
	})

	successfulCustomerRegistrations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "successful_customer_registration_count",
		Help: "The number of successful customer registration attempts since process start",
	})

	prometheus.MustRegister(lastResetTimestamp)
	prometheus.MustRegister(failedCustomerRegistrations)
	prometheus.MustRegister(successfulCustomerRegistrations)
}

// App defines shared dependencies, request handlers,
// and url-mappings for the API
type App struct {
	session      *mgo.Session
	redis        *redis.Client
	logger       *zap.Logger
	upgrader     *websocket.Upgrader
	mu           sync.Mutex
	counters     map[string]*websocket.Conn
	voiceSrvConn *websocket.Conn
	nextTicket   int
	jwtIssuer    string
	jwtSecret    string
}

// Router returns a http.Handler with url mappings
// for all routes handlers in the API
func (a *App) Router() http.Handler {
	r := chi.NewRouter()

	r.Post("/auth", a.login())

	r.Get("/users", a.listUsers())
	r.Post("/users", a.createUser())

	r.Get("/queues", a.listQueues())
	r.Post("/queues", a.createQueue())

	r.Get("/customers", a.listCustomers())
	r.Post("/customers", a.createCustomer())

	r.Post("/actions/next", a.callNextCustomer())

	// r.Get("/kiosks/new", a.handleKioskSetup)
	// r.Get("/counters/new", a.handleCounterSetup)

	r.Get("/ws", a.handleCounterWebsocket)
	r.Get("/voicews", a.handleVoiceWebsocket)

	r.Handle("/metrics", promhttp.Handler())

	fileServer(r, "/static", http.Dir("./ui/static"))

	return r
}

// ResetTicketing resets nextTicket
func (a *App) ResetTicketing() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nextTicket = 0
	lastResetTimestamp.Set(float64(time.Now().Unix()))
}

// NewApp returns a pointer to a new app with session and logger
// and websocket upgrader properly configured and ready for use in all routes
func NewApp(s *mgo.Session, c *redis.Client, l *zap.Logger, u *websocket.Upgrader, issuer, secret string) *App {
	a := App{
		session:    s,
		redis:      c,
		logger:     l,
		upgrader:   u,
		nextTicket: 0,
		jwtIssuer:  issuer,
		jwtSecret:  secret,
	}

	return &a
}

func (a *App) addCounter(id string, conn *websocket.Conn) {
	if a.counters == nil {
		a.counters = make(map[string]*websocket.Conn)
	}
	a.counters[id] = conn
}

func (a *App) getCounterConn(id string) *websocket.Conn {
	return a.counters[id]
}
