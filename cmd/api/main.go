// Package main is the entry point of the Go Example REST API.
//
//	@title			Go Example API
//	@version		1.0
//	@description	A clean-architecture REST API built with Gin and PostgreSQL.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	Abdul Falaq
//	@contact.url	https://github.com/abdulfalaq5
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@schemes	http https
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/falaqmsi/go-example/docs" // swag generated docs

	"github.com/falaqmsi/go-example/internal/config"
	"github.com/falaqmsi/go-example/internal/handler"
	"github.com/falaqmsi/go-example/internal/middleware"
	"github.com/falaqmsi/go-example/internal/repository"
	"github.com/falaqmsi/go-example/internal/service"
	"github.com/falaqmsi/go-example/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/MicahParks/keyfunc/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ── Gin mode ──────────────────────────────────────────────────────────────
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ── Database connections ───────────────────────────────────────────────────
	ctx := context.Background()
	db, err := storage.Connect(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// ── Storage Backend ────────────────────────────────────────────────────────
	fileStore, err := storage.NewFileStorage(cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize file storage module: %v", err)
	}

	// ── Auth (Keycloak JWKS) ───────────────────────────────────────────────────
	jwks, err := keyfunc.NewDefault([]string{cfg.KeycloakJWKSURL})
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %v", err)
	}
	authMiddleware := middleware.Auth(jwks)

	// ── Wire dependencies ──────────────────────────────────────────────────────

	// Health
	healthRepo := repository.NewHealthRepository()
	healthSvc := service.NewHealthService(healthRepo)
	healthHandler := handler.NewHealthHandler(healthSvc, cfg.AppEnv)

	// User
	userRepo := repository.NewUserRepository(db.Main)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Upload
	uploadSvc := service.NewUploadService(fileStore)
	uploadHandler := handler.NewUploadHandler(uploadSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.Default())

	// ── Routes ────────────────────────────────────────────────────────────────

	// Swagger UI  →  GET /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// System
	r.GET("/health", healthHandler.Check)

	// API v1
	v1 := r.Group("/api/v1")
	userHandler.RegisterRoutes(v1, authMiddleware)
	uploadHandler.RegisterRoutes(v1, authMiddleware)

	// Static routes for local storage uploads fallback
	r.Static("/uploads", cfg.Storage.LocalUploadDir)

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("route %s %s not found", c.Request.Method, c.Request.URL.Path),
		})
	})

	// ── HTTP server with graceful shutdown ────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
		log.Printf("📖 Swagger UI → http://localhost:%s/swagger/index.html", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server…")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Server exited cleanly.")
}
