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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddMetrics adds or updates metrics for a piece of content
func AddMetrics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized")
		return
	}

	userIDStr := userID.(string)
	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	var input models.MetricsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	contentObjID, err := primitive.ObjectIDFromHex(input.ContentID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid content ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify content belongs to user
	contentCollection := database.DB.Collection("generated_content")
	var content models.GeneratedContent
	err = contentCollection.FindOne(ctx, bson.M{
		"_id":     contentObjID,
		"user_id": userObjID,
	}).Decode(&content)

	if err == mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 404, "Content not found")
		return
	} else if err != nil {
		utils.ErrorResponse(c, 500, "Database error")
		return
	}

	// Calculate engagement rate
	totalEngagement := input.Likes + input.Comments + input.Shares
	engagementRate := 0.0
	if input.Impressions > 0 {
		engagementRate = (float64(totalEngagement) / float64(input.Impressions)) * 100
	}

	// Save metrics
	now := time.Now()
	metrics := models.ContentMetrics{
		ID:             primitive.NewObjectID(),
		ContentID:      contentObjID,
		UserID:         userObjID,
		Platform:       content.Platform,
		Likes:          input.Likes,
		Comments:       input.Comments,
		Shares:         input.Shares,
		Impressions:    input.Impressions,
		Reach:          input.Reach,
		EngagementRate: engagementRate,
		CollectedAt:    now,
		CreatedAt:      now,
	}

	metricsCollection := database.DB.Collection("content_metrics")
	_, err = metricsCollection.InsertOne(ctx, metrics)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to save metrics")
		return
	}

	utils.SuccessResponse(c, 201, "Metrics added successfully", metrics)
}

// GetMetricsByContentID retrieves metrics for a specific piece of content
func GetMetricsByContentID(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	contentID := c.Param("content_id")
	contentObjID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid content ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metricsCollection := database.DB.Collection("content_metrics")
	opts := options.Find().SetSort(bson.D{{Key: "collected_at", Value: -1}})

	cursor, err := metricsCollection.Find(ctx, bson.M{
		"content_id": contentObjID,
		"user_id":    userObjID,
	}, opts)

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch metrics")
		return
	}
	defer cursor.Close(ctx)

	var metrics []models.ContentMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode metrics")
		return
	}

	utils.SuccessResponse(c, 200, "Metrics retrieved successfully", metrics)
}

// GetAggregatedMetrics retrieves aggregated metrics for all user's content
func GetAggregatedMetrics(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metricsCollection := database.DB.Collection("content_metrics")

	// Aggregate metrics by platform
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userObjID}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$platform",
			"total_likes": bson.M{"$sum": "$likes"},
			"total_comments": bson.M{"$sum": "$comments"},
			"total_shares": bson.M{"$sum": "$shares"},
			"total_impressions": bson.M{"$sum": "$impressions"},
			"total_reach": bson.M{"$sum": "$reach"},
			"avg_engagement_rate": bson.M{"$avg": "$engagement_rate"},
			"post_count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := metricsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to aggregate metrics")
		return
	}
	defer cursor.Close(ctx)

	var platformMetrics []bson.M
	if err = cursor.All(ctx, &platformMetrics); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode aggregated metrics")
		return
	}

	// Calculate overall totals
	var totalLikes, totalComments, totalShares, totalImpressions, totalReach, totalPosts int
	var totalEngagementRate float64

	for _, pm := range platformMetrics {
		totalLikes += int(pm["total_likes"].(int32))
		totalComments += int(pm["total_comments"].(int32))
		totalShares += int(pm["total_shares"].(int32))
		totalImpressions += int(pm["total_impressions"].(int32))
		totalReach += int(pm["total_reach"].(int32))
		totalPosts += int(pm["post_count"].(int32))
		totalEngagementRate += pm["avg_engagement_rate"].(float64)
	}

	avgEngagementRate := 0.0
	if len(platformMetrics) > 0 {
		avgEngagementRate = totalEngagementRate / float64(len(platformMetrics))
	}

	response := map[string]interface{}{
		"total_posts":            totalPosts,
		"total_likes":            totalLikes,
		"total_comments":         totalComments,
		"total_shares":           totalShares,
		"total_impressions":      totalImpressions,
		"total_reach":            totalReach,
		"avg_engagement_rate":    avgEngagementRate,
		"by_platform":            platformMetrics,
	}

	utils.SuccessResponse(c, 200, "Aggregated metrics retrieved successfully", response)
}

