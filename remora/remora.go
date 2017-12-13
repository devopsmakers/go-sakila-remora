package remora

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Remora interface
type Remora struct {
	Config *Config
}

// Config basic remora configuration
type Config struct {
	Service       Connection
	AcceptableLag int
	CacheTTL      string
	HTTPServe     int
	Maintenance   bool
}

// Connection detail for MySQL
type Connection struct {
	Port int
	Host string
	User string
	Pass string
	Ssl  bool
}

// Result is the result of healthchecks
type Result struct {
	StatusCode int
	Body       bytes.Buffer
}

// HealthCheck interface that all "remoras" implement
type HealthCheck interface {
	Check() Result
}

type healthfunc func(*Config) Result

type statusHandler struct {
	r           *Remora
	healthfn    healthfunc
	cache       *cache.Cache
	mutex       *sync.Mutex
	contentType string
}

// LoadConfig inits the config file and reads the default config information
// into Remora.Config. For testability it accepts an array containing dirs to
// search for a config file.
func (r *Remora) LoadConfig(configpaths []string, servicename string) error {

	// Explicitly reset viper - helps for testing errors
	viper.Reset()

	viper.SetConfigName("config")

	for _, configpath := range configpaths {
		viper.AddConfigPath(configpath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.RegisterAlias(servicename, "Service")

	// some of our config names need aliasing to unmarshal correctly
	viper.RegisterAlias("acceptable-lag", "AcceptableLag")
	viper.RegisterAlias("cache-ttl", "CacheTTL")
	viper.RegisterAlias("http-serve", "HTTPServe")
	viper.RegisterAlias("maintenance-mode", "Maintenance")

	viper.SetDefault("Maintenance", false)
	viper.SetDefault("CacheTTL", "5s")

	if err := viper.Unmarshal(&r.Config); err != nil {
		return err
	}

	return nil
}

// Serve starts running checks and exposes the HTTP endpoint
func (r *Remora) Serve(healthfn healthfunc) error {

	// Lets create a cache with a ttl of the configured time
	ttl, err := time.ParseDuration(r.Config.CacheTTL)
	if err != nil {
		return err
	}

	cache := cache.New(ttl, 1*time.Second)

	status := statusHandler{
		r:        r,
		healthfn: healthfn,
		cache:    cache,
		mutex:    &sync.Mutex{},
	}

	listenPort := r.Config.HTTPServe

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", listenPort),
		Handler:        status,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	jww.WARN.Printf("starting remora listening on port: %d", listenPort)

	srv.SetKeepAlivesEnabled(false)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// ServeHTTP serves from cache, runs checks if cache empty
func (s statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jww.INFO.Printf("%s: requesting status", r.RemoteAddr)

	var resp Result
	tmp, found := s.cache.Get("res")

	if found {
		w.Header().Set("X-Cache", "HIT")
		resp = tmp.(Result)
	} else {

		s.mutex.Lock()
		defer s.mutex.Unlock()
		tmp, found := s.cache.Get("res")
		if found {
			w.Header().Set("X-Cache", "HIT")
			resp = tmp.(Result)
		} else {

			jww.INFO.Printf("%v: not in cache, getting status", r.RemoteAddr)
			w.Header().Set("X-Cache", "MISS")
			resp = s.healthfn(s.r.Config)

			s.cache.Set("res", resp, cache.DefaultExpiration)

		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Server", "remora-mysql")

	if resp.StatusCode == 0 {
		w.WriteHeader(http.StatusOK)
		resp.Body.WriteTo(w)
	} else if resp.StatusCode == 1 {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.Body.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		resp.Body.WriteTo(w)
	}

}
