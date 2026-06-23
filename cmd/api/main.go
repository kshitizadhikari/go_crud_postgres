package main

import (
	"go_crud_postgres/internal/config"
	"go_crud_postgres/internal/database"
	"go_crud_postgres/internal/routes"
	"go_crud_postgres/internal/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig()
	err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to datbase:", err)
	}

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	routes.RegisterRoutes(router, userHandler)
	log.Printf("🚀 Server running on http://localhost:%s", config.AppPort)
	if err := router.Run(":" + config.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
