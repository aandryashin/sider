package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	listen      string
	gracePeriod time.Duration
)

type handlerMethods map[string]http.Handler

func allowed(m handlerMethods) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, ok := m[r.Method]
		if !ok {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func clientDisconnectHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			handler.ServeHTTP(w, r)
			return
		}
		cnch := cn.CloseNotify()
		done := make(chan struct{})
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		go func() {
			handler.ServeHTTP(w, r.WithContext(ctx))
			done <- struct{}{}
		}()
		select {
		case <-cnch:
			cancel()
		case <-done:
		}
	})
}

func withParams(fn func(http.ResponseWriter, *http.Request, string, time.Duration)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.Split(r.URL.Path, "/")[2]
		if key == "" {
			http.Error(w, "Empty key.", http.StatusBadRequest)
			return
		}
		var ttl time.Duration
		var err error
		ttlStr := r.FormValue("ttl")
		if ttlStr != "" {
			ttl, err = time.ParseDuration(ttlStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Parse TTL: %v", err), http.StatusBadRequest)
				return
			}
			if ttl <= 0 {
				http.Error(w, "Zero or negative TTL.", http.StatusBadRequest)
				return
			}
		}
		fn(w, r, key, ttl)
	})
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/keys", allowed(
		handlerMethods{
			http.MethodGet: http.HandlerFunc(list),
		}))
	mux.Handle("/keys/", allowed(
		handlerMethods{
			http.MethodGet:    withParams(get),
			http.MethodPost:   withParams(set),
			http.MethodDelete: withParams(del),
		}))

	root := http.NewServeMux()
	root.Handle("/", clientDisconnectHandler(mux))

	return root
}

func init() {
	flag.StringVar(&listen, "listen", ":8080", "address to listel on")
	flag.DurationVar(&gracePeriod, "grace-period", 30*time.Second, "graceful shutdown period")
	flag.Parse()
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL)

	server := &http.Server{Addr: listen, Handler: handler()}
	go server.ListenAndServe()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
