package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo/internal/handler"
	"todo/internal/repository"
	"todo/internal/router"
	"todo/internal/service"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required")
	}

	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authH := handler.NewAuthHandler(authSvc)

	r := router.New(h, authH)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
