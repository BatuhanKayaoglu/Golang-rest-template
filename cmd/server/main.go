package main

import (
	"context"
	"golang-rest-api-template/pkg/api"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/config"
	"golang-rest-api-template/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.Load()

	redisClient := cache.NewRedisClient()
	db := database.NewDatabase()
	dbWrapper := &database.GormDatabase{DB: db}
	mongo := database.SetupMongoDB()
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := api.NewRouter(logger, mongo, dbWrapper, redisClient, &ctx)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
