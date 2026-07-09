package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go_crud_postgres/internal/config"
	"go_crud_postgres/internal/database"
	"go_crud_postgres/internal/routes"
	"go_crud_postgres/internal/user"
	"go_crud_postgres/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to datbase:", err)
	}

	ctx := context.Background()
	minioStorage, err := storage.NewMinIO(ctx, storage.Config{
		Endpoint:  cfg.Minio.Endpoint,
		AccessKey: cfg.Minio.AccessKey,
		SecretKey: cfg.Minio.SecretKey,
		UseSSL:    cfg.Minio.UseSSL,
		Bucket:    cfg.Minio.Bucket,
	})

	fmt.Println(minioStorage)
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository, minioStorage)
	userHandler := user.NewUserHandler(userService)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	routes.RegisterRoutes(router, userHandler)
	log.Printf("🚀 Server running on http://localhost:%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
