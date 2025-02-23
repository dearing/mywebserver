package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"text/template"
	"time"
)

var argVersion = flag.Bool("version", false, "print version information")

var bind = flag.String("bind", ":8080", "bind address")
var sseDuration = flag.Duration("sse-duration", 1*time.Second, "SSE ticker duration")

//go:embed embeded
var embedFS embed.FS

func usage() {
	println(`Usage: mywebserver [options]

Demo portable webserver with server side events and embedded assets.

- https://github.com/dearing/mywebserver

Options:
`)
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	// if the version flag is set, print version information and exit
	if *argVersion {
		versionInfo()
		return
	}

	// get a subtree fs of our embedded fs at wwwFS for static hosting
	wwwFS, err := fs.Sub(embedFS, "embeded/wwwroot")
	if err != nil {
		slog.Error("main/embedfs/wwwroot", "error", err)
		return
	}

	// get a subtree fs of our embedded fs at template for templates
	templateFS, err := fs.Sub(embedFS, "embeded/template")
	if err != nil {
		slog.Error("main/embedfs/template", "error", err)
		return
	}

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
		fmt.Fprintf(w, "hello world!")
	})

	// report embedded debug information about ourselves via a template
	handler.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("version called")

		// templatefs should be a subtree of our embedded fs
		template, err := template.ParseFS(templateFS, "version.html")
		if err != nil {
			slog.Error("version/template parse", "error", err)
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}

		// at build, go *can* embded handy information about the build
		info, ok := debug.ReadBuildInfo()
		if !ok {
			slog.Error("version/build info", "error", err)
			http.Error(w, "no build info", http.StatusInternalServerError)
			return
		}

		// with our info object ready, we can toss it over to the template to render
		err = template.Execute(w, info)
		if err != nil {
			slog.Error("version/template execute", "error", err)
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}

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

	// server static files from our subtree fs of wwwroot
	handler.Handle("/", http.FileServerFS(wwwFS))
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

	// we block until a signal is received
	<-sigchan

	// the server gets 2 seconds to shut itself down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown", "error", err)
	}

	slog.Info("server stopped")
}

func versionInfo() {
	// seems like a nice place to sneak in some debug information
	info, ok := debug.ReadBuildInfo()
	if ok {
		slog.Info("build info", "main", info.Main.Path, "version", info.Main.Version)
		for _, setting := range info.Settings {
			slog.Info("build info", "key", setting.Key, "value", setting.Value)
		}
	}
}
