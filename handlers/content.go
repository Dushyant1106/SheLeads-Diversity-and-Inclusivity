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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenerateBlog generates a blog post using AI
func GenerateBlog(c *gin.Context) {
	// Get user ID from query parameters
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

	var request models.BlogGenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get business profile
	businessCollection := database.DB.Collection("business_profiles")
	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&businessProfile)
	if err == mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 404, "Business profile not found. Please create one first.")
		return
	} else if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch business profile")
		return
	}

	// Generate blog content using AI
	contentAI := services.NewContentAIService()
	if contentAI == nil {
		utils.ErrorResponse(c, 500, "Content generation service unavailable")
		return
	}

	title, content, err := contentAI.GenerateBlogPost(ctx, &businessProfile, &request)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate blog: "+err.Error())
		return
	}

	// Save generated content
	now := time.Now()
	generatedContent := models.GeneratedContent{
		ID:         primitive.NewObjectID(),
		UserID:     userObjID,
		BusinessID: businessProfile.ID,
		Type:       models.ContentTypeBlog,
		Status:     models.ContentStatusGenerated,
		Title:      title,
		Content:    content,
		Metadata: map[string]interface{}{
			"topic":    request.Topic,
			"keywords": request.Keywords,
			"tone":     request.Tone,
			"length":   request.Length,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	contentCollection := database.DB.Collection("generated_content")
	_, err = contentCollection.InsertOne(ctx, generatedContent)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to save generated content")
		return
	}

	utils.SuccessResponse(c, 201, "Blog generated successfully", generatedContent)
}

// GenerateSocialPost generates social media posts using AI
func GenerateSocialPost(c *gin.Context) {
	// Get user ID from query parameters
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

	var request models.SocialPostGenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get business profile
	businessCollection := database.DB.Collection("business_profiles")
	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&businessProfile)
	if err == mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 404, "Business profile not found. Please create one first.")
		return
	} else if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch business profile")
		return
	}

	// Generate social posts using AI
	contentAI := services.NewContentAIService()
	if contentAI == nil {
		utils.ErrorResponse(c, 500, "Content generation service unavailable")
		return
	}

	posts, err := contentAI.GenerateSocialPost(ctx, &businessProfile, &request)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate social posts: "+err.Error())
		return
	}

	// Save generated posts
	contentCollection := database.DB.Collection("generated_content")
	var savedPosts []models.GeneratedContent
	now := time.Now()

	for platform, postContent := range posts {
		generatedContent := models.GeneratedContent{
			ID:         primitive.NewObjectID(),
			UserID:     userObjID,
			BusinessID: businessProfile.ID,
			Type:       models.ContentTypeSocial,
			Status:     models.ContentStatusGenerated,
			Title:      request.Topic,
			Content:    postContent,
			Platform:   platform,
			Hashtags:   request.Hashtags,
			Metadata: map[string]interface{}{
				"topic":       request.Topic,
				"tone":        request.Tone,
				"image_style": request.ImageStyle,
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = contentCollection.InsertOne(ctx, generatedContent)
		if err == nil {
			savedPosts = append(savedPosts, generatedContent)
		}
	}

	utils.SuccessResponse(c, 201, "Social posts generated successfully", savedPosts)
}

// GetAllContent retrieves all generated content for the authenticated user
func GetAllContent(c *gin.Context) {
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

	// Get query parameters for filtering
	contentType := c.Query("type")     // blog or social
	platform := c.Query("platform")    // linkedin, twitter, etc.
	status := c.Query("status")        // draft, generated, posted

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{"user_id": userObjID}
	if contentType != "" {
		filter["type"] = contentType
	}
	if platform != "" {
		filter["platform"] = platform
	}
	if status != "" {
		filter["status"] = status
	}

	contentCollection := database.DB.Collection("generated_content")
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := contentCollection.Find(ctx, filter, opts)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch content")
		return
	}
	defer cursor.Close(ctx)

	var contents []models.GeneratedContent
	if err = cursor.All(ctx, &contents); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode content")
		return
	}

	utils.SuccessResponse(c, 200, "Content retrieved successfully", contents)
}

// GetContentByID retrieves a specific piece of content by ID
func GetContentByID(c *gin.Context) {
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

	contentID := c.Param("id")
	contentObjID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid content ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	utils.SuccessResponse(c, 200, "Content retrieved successfully", content)
}

// DeleteContent deletes a specific piece of content
func DeleteContent(c *gin.Context) {
	// Get user ID from query parameters
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

	contentID := c.Param("id")
	contentObjID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid content ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contentCollection := database.DB.Collection("generated_content")
	result, err := contentCollection.DeleteOne(ctx, bson.M{
		"_id":     contentObjID,
		"user_id": userObjID,
	})

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete content")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, 404, "Content not found")
		return
	}

	utils.SuccessResponse(c, 200, "Content deleted successfully", nil)
}

// UpdateContentStatus updates the status of content (e.g., mark as posted)
func UpdateContentStatus(c *gin.Context) {
	// Get user ID from query parameters
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

	contentID := c.Param("id")
	contentObjID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid content ID")
		return
	}

	var input struct {
		Status  string `json:"status" binding:"required"`
		PostURL string `json:"post_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contentCollection := database.DB.Collection("generated_content")

	update := bson.M{
		"$set": bson.M{
			"status":     input.Status,
			"updated_at": time.Now(),
		},
	}

	if input.PostURL != "" {
		update["$set"].(bson.M)["post_url"] = input.PostURL
	}

	if input.Status == "posted" {
		now := time.Now()
		update["$set"].(bson.M)["posted_at"] = now
	}

	result, err := contentCollection.UpdateOne(ctx, bson.M{
		"_id":     contentObjID,
		"user_id": userObjID,
	}, update)

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update content status")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, 404, "Content not found")
		return
	}

	utils.SuccessResponse(c, 200, "Content status updated successfully", nil)
}

