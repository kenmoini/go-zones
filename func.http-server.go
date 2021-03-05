package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewRouter generates the router used in the HTTP Server
func NewRouter(basePath string) *http.ServeMux {
	if basePath == "" {
		basePath = "/" + appName
	}
	// Create router and define routes and return that router
	router := http.NewServeMux()

	// Version Output - reads from variables.go
	router.HandleFunc(basePath+"/version", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/version", r)
		fmt.Fprintf(w, appName+" version: %s\n", appVersion)
	})

	// Healthz endpoint for kubernetes platforms
	router.HandleFunc(basePath+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/healthz", r)
		fmt.Fprintf(w, "OK")
	})

	return router
}

// RunHTTPServer will run the HTTP Server
func (config Config) RunHTTPServer() {
	// Set up a channel to listen to for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Application.Server.Timeout.Server,
	)
	defer cancel()

	// Define server options
	server := &http.Server{
		Addr:         config.Application.Server.Host + ":" + config.Application.Server.Port,
		Handler:      NewRouter(config.Application.Server.BasePath),
		ReadTimeout:  config.Application.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Application.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Application.Server.Timeout.Idle * time.Second,
	}

	// Only listen on IPV4
	l, err := net.Listen("tcp4", config.Application.Server.Host+":"+config.Application.Server.Port)
	check(err)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Alert the user that the server is starting
	log.Printf("Server is starting on %s\n", server.Addr)

	// Run the server on a new goroutine
	go func() {
		//if err := server.ListenAndServe(); err != nil {
		if err := server.Serve(l); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
		}
	}()

	// Block on this channel listeninf for those previously defined syscalls assign
	// to variable so we can let the user know why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
