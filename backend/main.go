package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"thirdtask/backend/handlers"
	"thirdtask/backend/middleware"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true}))
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/items", handlers.GetItems)
		api.POST("/items", handlers.CreateItem)
		api.PUT("/items/:id", handlers.UpdateItem)
		api.DELETE("/items/:id", handlers.DeleteItem)
	}
	r.Run(":8080")
}
