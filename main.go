package main

import (
	"Golang/config"
	"Golang/controllers"
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
		api.POST("/events", controllers.CreateEvent)
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventById)
		api.PUT("events/:id", controllers.UpdateEvent)
		api.DELETE("events/:id", controllers.DeleteEvent)
	}

	server.Run(":8080")
}
