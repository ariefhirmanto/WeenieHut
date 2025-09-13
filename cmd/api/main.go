package main

import (
	"WeenieHut/internal/database"
	imagecompressor "WeenieHut/internal/image_compressor"
	"WeenieHut/internal/repository"
	"WeenieHut/internal/service"
	"WeenieHut/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"WeenieHut/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	db := database.New()
	defer db.Close()
	repo := repository.New(db)
	storage := storage.New(storage.S3Endpoint, storage.S3AccessKeyID, storage.S3SecretAccessKey, storage.Option{MaxConcurrent: 5})
	imageCompressor := imagecompressor.New(5)
	svc := service.New(repo, storage, imageCompressor)
	serv := server.NewServer(svc)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(serv, done)

	err := serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
