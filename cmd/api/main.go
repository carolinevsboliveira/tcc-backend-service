package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/caroline-oliveira/backend-tcc/internal/adapters/http"
	"github.com/caroline-oliveira/backend-tcc/internal/adapters/jwt"
	"github.com/caroline-oliveira/backend-tcc/internal/adapters/memory"
	"github.com/caroline-oliveira/backend-tcc/internal/application/auth"
	"github.com/caroline-oliveira/backend-tcc/internal/application/tasks"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	_ = godotenv.Load()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	seedPassword := os.Getenv("SEED_PASSWORD")
	if seedPassword == "" {
		log.Fatal("SEED_PASSWORD is required")
	}

	userRepository, err := memory.NewUserRepository(seedPassword)
	if err != nil {
		log.Fatal(err)
	}
	taskRepository := memory.NewTaskRepository()

	jwtProvider, err := jwt.NewJWTProvider([]byte(os.Getenv("JWT_SECRET")), time.Now)
	if err != nil {
		log.Fatal(err)
	}
	authService := auth.NewAuthService(userRepository, jwtProvider)
	taskService := tasks.NewService(taskRepository, userRepository, authService)
	router := httpadapter.NewRouter(httpadapter.Dependencies{
		Auth:  httpadapter.NewAuthHandler(authService),
		Tasks: httpadapter.NewTasksHandler(taskService),
	})
	server := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("backend-tcc listening on http://127.0.0.1%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
