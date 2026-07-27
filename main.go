package main

import (
	"Golang/config"
	"Golang/controllers"
	"Golang/middlewares"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()

	// Route Group
	api := server.Group("/api")
	{
		// Route Event
		api.POST("/events", controllers.CreateEvent)
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventById)
		api.PUT("/events/:id", controllers.UpdateEvent)
		api.DELETE("/events/:id", controllers.DeleteEvent)

		// Route Auth
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.RequireAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)
		}
	}

	server.Run(":8080")
}
