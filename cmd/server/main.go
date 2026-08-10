package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"decode-service/internal/config"
	"decode-service/internal/handlers"
	"decode-service/internal/services"
	"decode-service/api/v1"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()
	
	// Initialize services
	cryptoService := services.NewCryptoService()
	validationService := services.NewValidationService()
	
	// Initialize handlers
	decryptHandler := handlers.NewDecryptHandler(cryptoService, validationService)
	encryptHandler := handlers.NewEncryptHandler(cryptoService, validationService)
	healthHandler := handlers.NewHealthHandler()
	
	// Setup routes
	router := v1.SetupRoutes(cfg, decryptHandler, encryptHandler, healthHandler)
	
	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("API documentation available at: http://localhost:%s/swagger.html", cfg.Server.Port)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Server.Port)
	
	// Graceful shutdown
	go func() {
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Server is shutting down...")
}