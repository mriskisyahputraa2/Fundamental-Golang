package controllers

import (
	"Golang/config"
	"Golang/models"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

// Function ImageKit Client
func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

// Function untuk membuat event
func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Menerima File dari form data
	file, header, err := c.Request.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File Image is Required",
		})
		return
	}

	defer file.Close()

	// 1. Upload file ke imagekit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to Upload Image",
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	// 2. Simpan ke Database
	event := models.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserID:      userID.(int),
		Datetime:    parsedTime,
	}

	config.DB.Create(&event)
	c.JSON(http.StatusOK, gin.H{
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
	userID, _ := context.Get("userID")
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event Not Found",
		})
		return
	}

	if event.UserID != userID.(int) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to update this event",
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
	useID, _ := context.Get("userID")
	var event models.Event
	paramsId := context.Param("id")

	var EventData = config.DB.First(&event, paramsId).Error
	if EventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event Not Found",
		})
		return
	}

	if event.UserID != useID.(int) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this event",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
		"event":   event,
	})
}
