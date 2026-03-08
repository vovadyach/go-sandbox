package main

import (
	"fmt"
	"go-sandbox/api/internal/post"
	"go-sandbox/api/internal/response"
	"go-sandbox/api/internal/user"
	"log"
	"net/http"
	"time"

	"go-sandbox/api/internal/database"

	"go-sandbox/api/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// Repositories
	userRepo := user.NewRepository(db)
	postRepo := post.NewRepository(db)

	// Handlers
	userHandler := user.NewHandler(userRepo, cfg)
	postHandler := post.NewHandler(postRepo, cfg)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/users", userHandler.List)
		r.Get("/users/{id}", userHandler.GetByID)
		r.Get("/users/{id}/posts", postHandler.ListByUserID)
		r.Get("/posts", postHandler.List)
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Server is running on http://localhost%s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
