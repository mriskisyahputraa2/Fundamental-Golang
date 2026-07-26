package controllers

import (
	"Golang/config"
	"Golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function untuk membuat event
func CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	event.UserID = 1
	config.DB.Create(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// Function untuk mengambil data event
func GetEvents(context *gin.Context) {
	var events []models.Event
	config.DB.Find(&events)

	context.JSON(http.StatusOK, gin.H{
		"message": "Event retrieved successfully",
		"event":   events,
	})
}

// Function untuk mengambil data event berdasarkan ID
func GetEventById(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event Not Found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event Details Successfully",
		"event":   event,
	})

}

// Function Update Event
func UpdateEvent(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event Not Found",
		})
		return
	}

	var input models.Event
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(&input)

	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   event,
	})
}

// Function Delete Event
func DeleteEvent(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var EventData = config.DB.First(&event, paramsId).Error
	if EventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event Not Found",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
		"event":   event,
	})
}
