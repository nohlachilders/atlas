package server

import (
	"fmt"
	"net/http"
	"time"
)

type AddMiddlewareFunc func(http.Handler) http.Handler

func ChainMiddlewares(last http.Handler, middlewares []AddMiddlewareFunc) http.Handler {
	handler := last
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func AddLoggingMiddleware(next http.Handler) http.Handler {
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
		fmt.Printf("Request received: %v, Took: %v\n", start, time.Since(start))
	}()

	handler.next.ServeHTTP(w, r)
}
