package main

import (
	"Gravitum/internal/controller"
	"Gravitum/internal/database"
	"Gravitum/internal/repo"
	"Gravitum/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := gin.Default()
	api := router.Group("/api")
	userController.RegisterUserRoutes(api)

	port := ":8080"
	fmt.Printf("Server is running on %s...\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
