package server

import (
	"fmt"
	"net/http"
	"time"
)

// AddMiddlewareFunc is a function that adds middleware to an http.Handler
// functions of this type are intended to act as constructors for the associated
// middleware.
type AddMiddlewareFunc func(http.Handler, *Config) http.Handler

// ChainMiddlewares chains multiple middlewares together
// as suggested by the argument name, "last" is the last handler called, and
// is usually the endpoint being wrapped. the end of the list of given wrappers
// is called first.
func ChainMiddlewares(last http.Handler, middlewares []AddMiddlewareFunc, cfg *Config) http.Handler {
	handler := last
	for _, middleware := range middlewares {
		handler = middleware(handler, cfg)
	}
	return handler
}

func AddLoggingMiddleware(next http.Handler, cfg *Config) http.Handler {
	return &LoggingMiddleware{
		next: next,
	}
}

type LoggingMiddleware struct {
	next http.Handler
}

func (handler *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		fmt.Printf("Request received at %v: %v, Took: %v\n", r.URL, start, time.Since(start))
	}()

	handler.next.ServeHTTP(w, r)
}
