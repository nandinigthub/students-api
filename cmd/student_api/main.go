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
	"github.com/nandinigthub/students-api/internal/http/handlers/student"
	"github.com/nandinigthub/students-api/internal/storage/sqlite"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// database setup
	// dependency injection - Dependency Injection is a design pattern, that helps you to decouple the external logic of your implementation.
	//It’s common an implementation needs an external API, or a database, etc.
	store, err := sqlite.New(cfg) // to switch db just change here and create a new <newdb> folder and implement funct that it requires like the syntax and all
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Setup router
	r := http.NewServeMux()

	// Define a handler for the root endpoint
	r.HandleFunc("GET /api/student", student.Home())                                  // http://localhost:8080/api/student
	r.HandleFunc("POST /api/student/create", student.New(store))                      // http://localhost:8080/api/student/create
	r.HandleFunc("GET /api/student/{id}", student.GetstudentbyId(store))              // http://localhost:8080/api/student/{id}
	r.HandleFunc("GET /api/student/all", student.GetallStudent(store))                // http://localhost:8080/api/student/all
	r.HandleFunc("DELETE /api/student/delete/{id}", student.DeletestudentbyId(store)) // http://localhost:8080/api/student/delete

	// Setup server

	server := http.Server{
		Addr:    cfg.Addr, // Get the server address from the config
		Handler: r,
	}

	// Log server start message
	// fmt.Println(cfg.Addr)
	slog.Info("server started", slog.String("address", cfg.Addr))
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

	err = server.Shutdown(ctx) // takes too much time, it may hang infinitly therefore we use context package

	if err != nil {
		slog.Error("failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("shutdown server succesfully")
	// Start the server and log any error

}
