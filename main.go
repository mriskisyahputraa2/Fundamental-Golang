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
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventById)

		// Route Auth
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		// Route Middleware
		protected := api.Group("/")
		protected.Use(middlewares.RequireAuth())
		{
			protected.GET("/users/events", controllers.GetEventByUser)
			protected.GET("/auth/me", controllers.GetCurrentUser)
			protected.POST("/events", controllers.CreateEvent)
			protected.PUT("/events/:id", controllers.UpdateEvent)
			protected.DELETE("/events/:id", controllers.DeleteEvent)
		}
	}

	server.Run(":8080")
}
