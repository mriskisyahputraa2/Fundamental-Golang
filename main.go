package main

import (
	"Golang/config"
	"Golang/models"
	"log"
	"net/http"

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
		api.POST("/events",createEvent)
		api.GET("/events", getEvents)

	}

	server.Run(":8080")
}

// Function Handler 
func getEvents(context *gin.Context){
	events := models.GetAllEvent()

	context.JSON(http.StatusOK, events)	
}

func createEvent (context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)


	// jika terjadi error, kirimkan response error
	if err != nil {
		context.JSON(http.StatusBadRequest,gin.H{
			"message" : "Unable to parse request data",
			"error"   : err,
		})		
		return
	}
	// Dummy UserId
	event.UserId = 1


	// Save
	event.Save()

	context.JSON(http.StatusOK, gin.H{
		"message" : "Event created successfully",
		"event"   : event,
	})	
}