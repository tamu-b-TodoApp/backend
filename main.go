package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo/internal/handler"
	"todo/internal/middleware"
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

	authMiddleware := middleware.Auth()

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authH := handler.NewAuthHandler(authSvc, authMiddleware)

	companyRepo := repository.NewCompanyRepository(db)
	companySvc := service.NewCompanyService(companyRepo)
	companyH := handler.NewCompanyHandler(companySvc, authMiddleware)

	teamRepo := repository.NewTeamRepository(db)
	teamSvc := service.NewTeamService(teamRepo)
	teamH := handler.NewTeamHandler(teamSvc, authMiddleware)

	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	todoH := handler.NewTodoHandler(todoSvc, teamSvc, authMiddleware)

	healthH := handler.NewHealthHandler()

	r := router.New(authH, companyH, teamH, todoH, healthH)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
