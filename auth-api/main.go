package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	gommonlog "github.com/labstack/gommon/log"
)

var (
	ErrHttpGenericMessage = echo.NewHTTPError(http.StatusInternalServerError, "something went wrong, please try again later")
	ErrWrongCredentials   = echo.NewHTTPError(http.StatusUnauthorized, "username or password is invalid")
	jwtSecret             = "myfancysecret"
)

// ====== helpers env ======
func firstNonEmpty(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}
func normalizeBase(u string, def string) string {
	base := strings.TrimSpace(u)
	if base == "" {
		base = def
	}
	// si viene "host:puerto", añade esquema
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	return base
}

// ====== breaker simple ======
type brkState int

const (
	closed brkState = iota
	open
	halfOpen
)

type breaker struct {
	mu            sync.Mutex
	state         brkState
	failures      int
	lastChanged   time.Time
	maxFailures   int
	openInterval  time.Duration
	halfOpenLimit int
	halfCalls     int
}

func newBreaker(maxFailures int, openInterval time.Duration, halfOpenLimit int) *breaker {
	return &breaker{
		state:         closed,
		maxFailures:   maxFailures,
		openInterval:  openInterval,
		halfOpenLimit: halfOpenLimit,
		lastChanged:   time.Now(),
	}
}
func (b *breaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case closed:
		return true
	case open:
		if time.Since(b.lastChanged) >= b.openInterval {
			b.state = halfOpen
			b.halfCalls = 0
			return true
		}
		return false
	case halfOpen:
		if b.halfCalls < b.halfOpenLimit {
			b.halfCalls++
			return true
		}
		return false
	default:
		return true
	}
}
func (b *breaker) success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == halfOpen {
		b.state = closed
		b.failures = 0
		b.lastChanged = time.Now()
	} else if b.state == closed {
		b.failures = 0
	}
}
func (b *breaker) failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case closed:
		b.failures++
		if b.failures >= b.maxFailures {
			b.state = open
			b.lastChanged = time.Now()
		}
	case halfOpen:
		b.state = open
		b.lastChanged = time.Now()
	}
}

// ====== tipos de tu proyecto ======
// NOTA: NO redefinimos UserService aquí. Ya existe en user.go (mismo package).

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// puerto
	port := os.Getenv("AUTH_API_PORT")
	if port == "" {
		port = "8081"
	}
	hostport := ":" + port

	// secret
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" {
		jwtSecret = v
	}

	// base de users-api (acepta muchos alias)
	rawBase := firstNonEmpty(
		"USERS_API_URL",
		"USERS_SERVICE_URL",
		"USER_SERVICE_URL",
		"USER_API_URL",
		"USERS_URL",
		"USER_API_ADDR",
		"USERS_SERVICE_ADDR",
		"USERS_SERVICE_HOST",
		"USERS_API_ADDRESS", // compatibilidad con tu main.go original
	)
	userAPIBase := normalizeBase(rawBase, "http://users-api:8083")

	// Validar URL
	if _, err := url.ParseRequestURI(userAPIBase); err != nil {
		log.Fatalf("invalid USERS API base URL: %q (%v)", userAPIBase, err)
	}

	// UserService es el que ya tienes en user.go
	userService := UserService{
		Client:         http.DefaultClient,
		UserAPIAddress: userAPIBase,
		AllowedUserHashes: map[string]interface{}{
			"admin_admin": nil,
			"johnd_foo":   nil,
			"janed_ddd":   nil,
		},
	}

	e := echo.New()
	e.Logger.SetLevel(gommonlog.INFO)

	// Zipkin (si tienes initTracing en tu repo)
	if zipkinURL := os.Getenv("ZIPKIN_URL"); zipkinURL != "" {
		e.Logger.Infof("init tracing to Zipkit at %s", zipkinURL)
		if tracedMiddleware, tc, err := initTracing(zipkinURL); err == nil {
			e.Use(echo.WrapMiddleware(tracedMiddleware))
			// Solo si el client retornado es *http.Client lo usamos; si no, seguimos con DefaultClient
			if c, ok := any(tc).(*http.Client); ok && c != nil {
				userService.Client = c
			} else {
				e.Logger.Infof("zipkin client is not *http.Client; keeping default transport")
			}
		} else {
			e.Logger.Infof("Zipkin tracer init failed: %s", err.Error())
		}
	} else {
		e.Logger.Infof("Zipkin URL was not provided, tracing is not initialised")
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// health
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })
	e.GET("/version", func(c echo.Context) error { return c.String(http.StatusOK, "Auth API, written in Go\n") })

	// breaker (parametrizable)
	brk := newBreaker(
		envInt("CB_MAX_FAILURES", 3),
		envDurMs("CB_OPEN_MS", 10000),
		envInt("CB_HALF_OPEN_LIMIT", 2),
	)

	e.POST("/login", getLoginHandler(userService, brk))

	e.Logger.Infof("auth-api on %s; USERS_API=%s", hostport, userAPIBase)
	e.Logger.Fatal(e.Start(hostport))
}

func getLoginHandler(userService UserService, brk *breaker) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusNoContent)
		}

		if !brk.allow() {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "auth service temporarily unavailable")
		}

		req := LoginRequest{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			brk.failure()
			return ErrHttpGenericMessage
		}

		ctx := c.Request().Context()
		user, err := userService.Login(ctx, strings.TrimSpace(req.Username), req.Password)
		if err != nil {
			if err != ErrWrongCredentials {
				brk.failure()
				log.Printf("could not authorize user '%s': %v", req.Username, err)
				return ErrHttpGenericMessage
			}
			// credenciales malas NO abren el breaker
			return ErrWrongCredentials
		}
		brk.success()

		// JWT HS256
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["firstname"] = user.FirstName
		claims["lastname"] = user.LastName
		claims["role"] = user.Role
		claims["exp"] = time.Now().Add(72 * time.Hour).Unix()

		t, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return ErrHttpGenericMessage
		}

		// Devolvemos ambas claves por conveniencia
		return c.JSON(http.StatusOK, map[string]string{
			"accessToken": t,
			"token":       t,
		})
	}
}

// ====== env parsers para breaker ======
func envInt(key string, def int) int {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err == nil && n >= 0 {
		return n
	}
	return def
}

func envDurMs(key string, defMs int) time.Duration {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return time.Duration(defMs) * time.Millisecond
	}
	n, err := strconv.Atoi(s)
	if err == nil && n >= 0 {
		return time.Duration(n) * time.Millisecond
	}
	return time.Duration(defMs) * time.Millisecond
}
