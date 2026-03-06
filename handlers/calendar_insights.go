package handlers

import (
	"context"
	"fmt"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"sheleads-backend/services"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FreeTimeSlot struct {
	Day       string `json:"day"`
	TimeRange string `json:"time_range"`
	Duration  string `json:"duration"`
}

type WorkPattern struct {
	BusiestDay     string  `json:"busiest_day"`
	AvgHoursPerDay float64 `json:"avg_hours_per_day"`
	PeakTime       string  `json:"peak_time"`
	FreeDays       int     `json:"free_days"`
}

type BusinessSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeSlot    string `json:"time_slot"`
	Duration    string `json:"duration"`
}

type CalendarInsightsResponse struct {
	FreeTimeSlots       []FreeTimeSlot       `json:"free_time_slots"`
	WorkPattern         WorkPattern          `json:"work_pattern"`
	BusinessSuggestions []BusinessSuggestion `json:"business_suggestions"`
}

// GetCalendarInsights analyzes work patterns and suggests free time
func GetCalendarInsights(c *gin.Context) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get work logs from the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	worklogsCollection := database.DB.Collection("worklogs")

	cursor, err := worklogsCollection.Find(ctx, bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": thirtyDaysAgo,
		},
	}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch work logs")
		return
	}
	defer cursor.Close(ctx)

	var workLogs []models.WorkLog
	if err = cursor.All(ctx, &workLogs); err != nil {
		utils.ErrorResponse(c, 500, "Failed to parse work logs")
		return
	}

	// Use AI to analyze patterns and generate insights
	geminiService, err := services.NewGeminiService(ctx)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to initialize AI service")
		return
	}
	defer geminiService.Close()

	// Get business profile for personalized suggestions
	var businessProfile models.BusinessProfile
	businessCollection := database.DB.Collection("business_profiles")
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&businessProfile)
	hasBusinessProfile := err == nil

	// Generate AI insights
	insights, err := geminiService.AnalyzeWorkPatterns(ctx, workLogs, businessProfile, hasBusinessProfile)
	if err != nil {
		utils.ErrorResponse(c, 500, fmt.Sprintf("Failed to analyze patterns: %v", err))
		return
	}

	utils.SuccessResponse(c, 200, "Calendar insights generated", insights)
}

