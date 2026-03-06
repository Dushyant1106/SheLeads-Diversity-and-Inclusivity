package handlers

import (
	"context"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignupRequest struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	Age              int    `json:"age" binding:"required,min=1"`
	EmergencyContact string `json:"emergency_contact" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  models.User  `json:"user"`
}

func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := database.DB.Collection("users")

	// Check if user already exists
	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		utils.ErrorResponse(c, 409, "User with this email already exists")
		return
	}

	// Create new user
	user := models.User{
		ID:               primitive.NewObjectID(),
		Name:             req.Name,
		Email:            req.Email,
		Age:              req.Age,
		EmergencyContact: req.EmergencyContact,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		utils.ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// Save to database
	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create user")
		return
	}

	utils.SuccessResponse(c, 201, "User created successfully", AuthResponse{
		User:  user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := database.DB.Collection("users")

	// Find user by email
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		utils.ErrorResponse(c, 401, "Invalid email or password")
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		utils.ErrorResponse(c, 401, "Invalid email or password")
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", AuthResponse{
		User:  user,
	})
}

func GetProfile(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := database.DB.Collection("users")

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	utils.SuccessResponse(c, 200, "Profile retrieved successfully", user)
}

