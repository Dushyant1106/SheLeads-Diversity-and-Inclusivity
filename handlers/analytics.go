package handlers

import (
	"context"
	"sheleads-backend/database"
	
	"sheleads-backend/models"
	"sheleads-backend/services"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityStats struct {
	TotalWorkLogs     int                                    `json:"total_work_logs"`
	TotalHours        float64                                `json:"total_hours"`
	TotalPoints       int                                    `json:"total_points"`
	AverageConfidence float64                                `json:"average_confidence"`
	CategoryBreakdown map[models.WorkCategory]CategoryStats  `json:"category_breakdown"`
	WeeklyTrend       []WeeklyStats                          `json:"weekly_trend"`
}

type WeeklyStats struct {
	Week   string  `json:"week"`
	Hours  float64 `json:"hours"`
	Points int     `json:"points"`
	Count  int     `json:"count"`
}

type CategoryStats struct {
	Count  int     `json:"count"`
	Hours  float64 `json:"hours"`
	Points int     `json:"points"`
}

func GetActivityStats(c *gin.Context) {
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

	// Get query parameters for date range
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
	} else {
		// Default to last 30 days
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
	} else {
		endDate = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch work logs
	worklogsCollection := database.DB.Collection("worklogs")
	cursor, err := worklogsCollection.Find(ctx, bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
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

	// Calculate statistics
	stats := calculateActivityStats(workLogs)

	utils.SuccessResponse(c, 200, "Activity stats retrieved successfully", stats)
}

func calculateActivityStats(workLogs []models.WorkLog) ActivityStats {
	stats := ActivityStats{
		CategoryBreakdown: make(map[models.WorkCategory]CategoryStats),
		WeeklyTrend:       []WeeklyStats{},
	}

	weeklyMap := make(map[string]*WeeklyStats)
	var totalConfidence float64

	for _, log := range workLogs {
		stats.TotalWorkLogs++
		stats.TotalHours += log.EstimatedHours
		stats.TotalPoints += log.Points
		totalConfidence += log.ConfidenceScore

		// Category breakdown
		catStats := stats.CategoryBreakdown[log.Category]
		catStats.Count++
		catStats.Hours += log.EstimatedHours
		catStats.Points += log.Points
		stats.CategoryBreakdown[log.Category] = catStats

		// Weekly trend
		year, week := log.CreatedAt.ISOWeek()
		weekKey := time.Date(year, 0, (week-1)*7+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		
		if weeklyMap[weekKey] == nil {
			weeklyMap[weekKey] = &WeeklyStats{Week: weekKey}
		}
		weeklyMap[weekKey].Hours += log.EstimatedHours
		weeklyMap[weekKey].Points += log.Points
		weeklyMap[weekKey].Count++
	}

	// Calculate average confidence
	if stats.TotalWorkLogs > 0 {
		stats.AverageConfidence = totalConfidence / float64(stats.TotalWorkLogs)
	}

	// Convert weekly map to slice
	for _, weekStats := range weeklyMap {
		stats.WeeklyTrend = append(stats.WeeklyTrend, *weekStats)
	}

	return stats
}

// GetAnalyticsSummary returns a summary of analytics for the dashboard
func GetAnalyticsSummary(c *gin.Context) {
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

	// Fetch all work logs for the user (last 30 days by default)
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	worklogsCollection := database.DB.Collection("worklogs")
	cursor, err := worklogsCollection.Find(ctx, bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
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

	// Calculate summary statistics
	var totalHours float64
	var totalPoints int
	categoryMap := make(map[string]map[string]interface{})

	for _, log := range workLogs {
		totalHours += log.EstimatedHours
		totalPoints += log.Points

		// Build category breakdown
		category := string(log.Category)
		if categoryMap[category] == nil {
			categoryMap[category] = map[string]interface{}{
				"category":    category,
				"total_hours": 0.0,
				"count":       0,
			}
		}
		categoryMap[category]["total_hours"] = categoryMap[category]["total_hours"].(float64) + log.EstimatedHours
		categoryMap[category]["count"] = categoryMap[category]["count"].(int) + 1
	}

	// Convert category map to array
	byCategory := []map[string]interface{}{}
	for _, cat := range categoryMap {
		byCategory = append(byCategory, cat)
	}

	summary := map[string]interface{}{
		"total_hours":  totalHours,
		"total_points": totalPoints,
		"by_category":  byCategory,
	}

	utils.SuccessResponse(c, 200, "Analytics summary retrieved successfully", summary)
}

func GetBurnoutStatus(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	burnoutService := services.NewBurnoutService()
	status, err := burnoutService.CheckBurnout(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to check burnout status")
		return
	}

	utils.SuccessResponse(c, 200, "Burnout status retrieved successfully", status)
}

