package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harungurubudi/mtsg/internal/di"
)

func main() {
	// Initialize the HTTP server using Wire
	server := di.InitializeServer()

	// Test the DI setup
	fmt.Println("✅ Echo server initialized successfully!")
	fmt.Printf("Server type: %T\n", server)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("🚀 MTSG Echo server started on port 8080")
	log.Println("📡 Ping endpoint: http://localhost:8080/ping")
	log.Println("🏥 Health endpoint: http://localhost:8080/health")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
