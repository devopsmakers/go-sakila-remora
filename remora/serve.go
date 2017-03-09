package remora

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	jww "github.com/spf13/jwalterweatherman"
)

type res struct {
	statusCode int
	body       bytes.Buffer
}

type statusHandler struct {
	r            *Remora
	remoraConfig *Config
	cache        *cache.Cache
	mutex        *sync.Mutex
	contentType  string
}

// Serve starts running checks and exposes the HTTP endpoint
func (r *Remora) Serve() error {

	// Lets create a cache with a ttl of the configured time
	ttl, err := time.ParseDuration(r.Config.CacheTTL)
	if err != nil {
		return err
	}

	cache := cache.New(ttl, 5*time.Second)

	status := statusHandler{
		r:            r,
		remoraConfig: r.Config,
		cache:        cache,
		mutex:        &sync.Mutex{},
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

	var resp res
	tmp, found := s.cache.Get("result")

	if found {
		resp = tmp.(res)
	} else {

		s.mutex.Lock()
		defer s.mutex.Unlock()
		tmp, found := s.cache.Get("res")
		if found {
			resp = tmp.(res)
		} else {

			jww.INFO.Printf("%v: stale cache, getting status", r.RemoteAddr)
			status, values, err := GetSlaveStatus(s.remoraConfig)
			if err != nil {
				jww.ERROR.Println(err)
			}

			var body = bytes.NewBuffer(values)
			resp = res{statusCode: status, body: *body}

			s.cache.Set("res", resp, cache.DefaultExpiration)

		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "go-sakila-remora")

	if resp.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
		resp.body.WriteTo(w)
	} else if resp.statusCode == 1 {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.body.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		resp.body.WriteTo(w)
	}

}
