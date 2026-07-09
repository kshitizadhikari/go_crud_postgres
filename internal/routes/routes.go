package routes

import (
	"go_crud_postgres/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userHandler *user.UserHandler) {
	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	users := v1.Group("/users")
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetById)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		// users.PATCH("/:id", userHandler.PatchUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
