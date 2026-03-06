package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sheleads-backend/database"
	"sheleads-backend/middleware"
	"sheleads-backend/models"
	"sheleads-backend/services"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogWorkRequest struct {
	Category    models.WorkCategory `form:"category" binding:"required"`
	Description string              `form:"description" binding:"required"`
}

func LogWork(c *gin.Context) {
	// Get user ID from form data
	userIDStr := c.PostForm("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	var req LogWorkRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, 400, "Image is required")
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		utils.ErrorResponse(c, 500, "Failed to create uploads directory")
		return
	}

	// Save image with unique filename
	filename := fmt.Sprintf("%s_%d_%s", userIDStr, time.Now().Unix(), filepath.Base(file.Filename))
	imagePath := filepath.Join(uploadsDir, filename)

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		utils.ErrorResponse(c, 500, "Failed to save image")
		return
	}

	// Verify work with Gemini AI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	geminiService, err := services.NewGeminiService(ctx)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to initialize AI service")
		return
	}
	defer geminiService.Close()

	verification, err := geminiService.VerifyWork(ctx, imagePath, req.Category, req.Description)
	if err != nil {
		// If AI verification fails, clean up the image and return error
		os.Remove(imagePath)
		utils.ErrorResponse(c, 500, fmt.Sprintf("AI verification failed: %v", err))
		return
	}

	// Create work log entry
	workLog := models.WorkLog{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Category:        req.Category,
		Description:     req.Description,
		ImagePath:       imagePath,
		ConfidenceScore: verification.ConfidenceScore,
		AIVerification:  verification.Explanation,
		EstimatedHours:  verification.EstimatedHours,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Calculate points
	workLog.CalculatePoints()

	// Save to database only if confidence score is acceptable
	if verification.ConfidenceScore < 0.3 {
		// Delete the image if verification failed
		os.Remove(imagePath)
		utils.ErrorResponse(c, 400, fmt.Sprintf("Work verification failed: %s", verification.Explanation))
		return
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	worklogsCollection := database.DB.Collection("worklogs")
	_, err = worklogsCollection.InsertOne(ctx2, workLog)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to save work log")
		return
	}

	// Check for burnout after logging work
	burnoutService := services.NewBurnoutService()
	burnoutStatus, _ := burnoutService.CheckBurnout(userIDStr)

	response := map[string]interface{}{
		"work_log":       workLog,
		"verification":   verification,
		"burnout_status": burnoutStatus,
	}

	utils.SuccessResponse(c, 201, "Work logged successfully", response)
}

func GetWorkLogs(c *gin.Context) {
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

	worklogsCollection := database.DB.Collection("worklogs")

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := worklogsCollection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch work logs")
		return
	}
	defer cursor.Close(ctx)

	var workLogs []models.WorkLog
	if err = cursor.All(ctx, &workLogs); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode work logs")
		return
	}

	if workLogs == nil {
		workLogs = []models.WorkLog{}
	}

	utils.SuccessResponse(c, 200, "Work logs retrieved successfully", workLogs)
}

func GetWorkLogByID(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	workLogIDStr := c.Param("id")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	workLogID, err := primitive.ObjectIDFromHex(workLogIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid work log ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	worklogsCollection := database.DB.Collection("worklogs")

	var workLog models.WorkLog
	err = worklogsCollection.FindOne(ctx, bson.M{"_id": workLogID, "user_id": userID}).Decode(&workLog)
	if err != nil {
		utils.ErrorResponse(c, 404, "Work log not found")
		return
	}

	utils.SuccessResponse(c, 200, "Work log retrieved successfully", workLog)
}

func DeleteWorkLog(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	workLogIDStr := c.Param("id")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	workLogID, err := primitive.ObjectIDFromHex(workLogIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid work log ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	worklogsCollection := database.DB.Collection("worklogs")

	var workLog models.WorkLog
	err = worklogsCollection.FindOne(ctx, bson.M{"_id": workLogID, "user_id": userID}).Decode(&workLog)
	if err != nil {
		utils.ErrorResponse(c, 404, "Work log not found")
		return
	}

	// Delete image file
	if workLog.ImagePath != "" {
		os.Remove(workLog.ImagePath)
	}

	// Delete from database
	_, err = worklogsCollection.DeleteOne(ctx, bson.M{"_id": workLogID, "user_id": userID})
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete work log")
		return
	}

	utils.SuccessResponse(c, 200, "Work log deleted successfully", nil)
}

