package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"time"

	"golang.org/x/net/context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	handleSig bool

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of all HTTP requests",
	}, []string{"code", "handler", "method"})
)

var buildVersion string

func root(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fmt.Fprintf(w, "<html><body><h1>Fleet Demo App</h1><p>Fleet demo app version %s</p></body></html>", buildVersion)
	log.Printf("Took %s to respond with buildVersion %s", time.Since(start), buildVersion)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if !handleSig {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	} else {
		w.WriteHeader(http.StatusGone)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
	}
}

func main() {
	r := mux.NewRouter()
	mainHandlerChain := promhttp.InstrumentHandlerDuration(
		httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "main"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, http.HandlerFunc(root)),
	)
	prom := prometheus.NewRegistry()
	prom.MustRegister(httpRequestsTotal)

	r.Handle("/", mainHandlerChain)
	r.HandleFunc("/health", healthCheck)
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	//Wait for sigterm
	<-c

	//Will cause healthcheck to unregister
	handleSig = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	log.Println("SIGTERM received.. waiting 10 seconds before shutdown")
	time.Sleep(10 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
