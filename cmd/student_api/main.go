package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nandinigthub/students-api/internal/config"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Setup router
	r := http.NewServeMux()

	// Define a handler for the root endpoint
	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home page of students API"))
	})

	// Setup server

	server := http.Server{
		Addr:    cfg.Addr, // Get the server address from the config
		Handler: r,
	}

	// Log server start message
	fmt.Println(cfg.Addr)
	fmt.Printf("Server started successfully, running at %s", cfg.HTTPServer.Addr)

	// graceful close/ shutdown
	done := make(chan os.Signal, 1)

	//a synchronous signal is converted into a run-time panic.
	// A SIGHUP, SIGINT, or SIGTERM signal causes the program to exit.
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	<-done

	slog.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx) // takes too much time, it may hang infinitly therefore we use context package

	if err != nil {
		slog.Error("failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("shutdown server succesfully")
	// Start the server and log any error

}
