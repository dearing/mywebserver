package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var bind = flag.String("bind", ":8080", "bind address")
var sseDuration = flag.Duration("sse-duration", 1*time.Second, "sse duration")

//go:embed wwwroot
var wwwroot embed.FS

func main() {

	flag.Parse()

	// create a channel to listen for signals
	sigchan := make(chan os.Signal, 1)
	// we want to capture SIGINT and SIGTERM signals and handle them gracefully
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// one pulse every duration
	ticker := time.NewTicker(*sseDuration)
	defer ticker.Stop()

	handler := http.NewServeMux()

	// simple hello handler
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("hello received")
		w.Write([]byte("hello world!"))
	})

	// server side events handler to distribute cars
	handler.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// check if the client supports flushing
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		// get the context for this request
		ctx := r.Context()

		slog.Info("client connected", "remote", r.RemoteAddr)

		// game loop to win a car
		for {
			select {
			// if the client disconnects, we should stop sending events
			case <-ctx.Done():
				slog.Info("client disconnected", "remote", r.RemoteAddr)
				return
			// on pulse, this client gets a car!
			case <-ticker.C:
				fmt.Fprintf(w, "data: %s gets a car!\n\n", r.RemoteAddr)
				flusher.Flush()
			}
		}
	})

	// mount our embedded wwwroot and serve
	handler.Handle("/", http.FileServerFS(wwwroot))
	server := &http.Server{
		Addr:    *bind,
		Handler: handler,
	}

	// spin off the the server in a goroutine, we can call shutdown on it later
	go func() {
		slog.Info("http server listening", "bind", *bind)
		server.ListenAndServe()
		slog.Info("http server stopped")
	}()

	<-sigchan

	// the server gets 10 seconds to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown", "error", err)
	}

	slog.Info("server stopped")
}
