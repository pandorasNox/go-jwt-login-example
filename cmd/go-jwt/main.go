package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// import "github.com/sirupsen/logrus"

func main() {
	addr := flag.String("addr", ":443", "address to bind to")
	dryRun := flag.Bool("dry-run", false, "enables dry-run mode, always returning success AdmissionReview")
	debug := flag.Bool("debug", false, "log all requests")
	flag.Parse()

	fmt.Println("starting...")
	log.WithFields(log.Fields{
		"debug":  *debug,
		"addr":   *addr,
		"dryRun": *dryRun,
	}).Info("programm flags")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/liveness", liveness)
	mux.HandleFunc("/", index)
	server := &http.Server{Addr: *addr, Handler: mux}

	go runServer(server)

	// wait for an exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Printf("waiting for server shutdown")
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
}

func runServer(server *http.Server) {
	log.Printf("server stopped: %s", server.ListenAndServe())
}

func healthz(w http.ResponseWriter, r *http.Request) {

	// if err != nil {
	// 	w.WriteHeader(http.StatusServiceUnavailable)
	// 	log.WithFields(log.Fields{"status": "unhealthy"}).Error("/healthz")
	// 	return
	// }

}

func liveness(w http.ResponseWriter, r *http.Request) {}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
