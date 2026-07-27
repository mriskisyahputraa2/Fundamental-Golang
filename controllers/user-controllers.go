package controllers

import (
	"Golang/config"
	"Golang/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister

	// Validation Register
	err := context.ShouldBindJSON(&input)
	if err != nil {
		// Parse pesan error per field
		var ve validator.ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			ve = errors
			for _, e := range ve {
				switch e.Field() {
				case "Email":
					context.JSON(http.StatusBadRequest, gin.H{
						"message": "Invalid email format",
					})
				case "Password":
					if e.Tag() == "min" {
						context.JSON(http.StatusBadRequest, gin.H{
							"message": "Password must be at least 6 characters",
						})
					} else {
						context.JSON(http.StatusBadRequest, gin.H{
							"message": "Password is required",
						})
					}
				default:
					context.JSON(http.StatusBadRequest, gin.H{
						"message": e.Field() + " is required",
					})
				}
				return
			}
		}
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hashing Password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Save User to Database
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error

	// Check if Email already exists
	if userCreated != nil {
		context.JSON(http.StatusConflict, gin.H{
			"message": "Email already taken",
		})
		return
	}

	// Message Success User Response
	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully registered",
		"user": gin.H{
			"name":     user.Name,
			"email":    user.Email,
			"datetime": user.CreatedAt,
		},
	})

}

func LoginUser(context *gin.Context) {
	var input AuthInputLogin

	err := context.ShouldBindJSON(&input)
	if err != nil {
		// Parse pesan error per field
		var ve validator.ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			ve = errors
			for _, e := range ve {
				switch e.Field() {
				case "Email":
					context.JSON(http.StatusBadRequest, gin.H{
						"message": "Invalid email format",
					})
				case "Password":
					if e.Tag() == "min" {
						context.JSON(http.StatusBadRequest, gin.H{
							"message": "Password must be at least 6 characters",
						})
					} else {
						context.JSON(http.StatusBadRequest, gin.H{
							"message": "Password is required",
						})
					}
				default:
					context.JSON(http.StatusBadRequest, gin.H{
						"message": e.Field() + " is required",
					})
				}
				return
			}
		}
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find User by Email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Check Password
	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchPassword != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT Token
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, errToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errToken != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Response Success
	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   signedToken,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})

}

// Get Profile
func GetCurrentUser(context *gin.Context) {

	// Get userId
	userId, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	var user models.User
	userData := config.DB.Preload("Events").First(&user, userId).Error

	if userData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "User Not Found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})

}
