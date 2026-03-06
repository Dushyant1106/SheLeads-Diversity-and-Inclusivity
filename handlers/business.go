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
)

// CreateBusinessProfile creates or updates a business profile for the authenticated user
func CreateBusinessProfile(c *gin.Context) {
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

	var input models.BusinessProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	businessCollection := database.DB.Collection("business_profiles")

	// Check if business profile already exists for this user
	var existingProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&existingProfile)

	now := time.Now()

	if err == mongo.ErrNoDocuments {
		// Create new business profile
		businessProfile := models.BusinessProfile{
			ID:                  primitive.NewObjectID(),
			UserID:              userObjID,
			BusinessName:        input.BusinessName,
			Industry:            input.Industry,
			Location:            input.Location,
			Description:         input.Description,
			TargetAudience:      input.TargetAudience,
			UniqueSellingPoints: input.UniqueSellingPoints,
			Website:             input.Website,
			SocialMediaHandles:  input.SocialMediaHandles,
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		_, err = businessCollection.InsertOne(ctx, businessProfile)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to create business profile")
			return
		}

		utils.SuccessResponse(c, 201, "Business profile created successfully", businessProfile)
	} else if err != nil {
		utils.ErrorResponse(c, 500, "Database error")
		return
	} else {
		// Update existing business profile
		update := bson.M{
			"$set": bson.M{
				"business_name":         input.BusinessName,
				"industry":              input.Industry,
				"location":              input.Location,
				"description":           input.Description,
				"target_audience":       input.TargetAudience,
				"unique_selling_points": input.UniqueSellingPoints,
				"website":               input.Website,
				"social_media_handles":  input.SocialMediaHandles,
				"updated_at":            now,
			},
		}

		_, err = businessCollection.UpdateOne(ctx, bson.M{"user_id": userObjID}, update)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to update business profile")
			return
		}

		// Fetch updated profile
		var updatedProfile models.BusinessProfile
		err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&updatedProfile)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to fetch updated profile")
			return
		}

		utils.SuccessResponse(c, 200, "Business profile updated successfully", updatedProfile)
	}
}

// GetBusinessProfile retrieves the business profile for the authenticated user
func GetBusinessProfile(c *gin.Context) {
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

	businessCollection := database.DB.Collection("business_profiles")

	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&businessProfile)

	if err == mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 404, "Business profile not found. Please create one first.")
		return
	} else if err != nil {
		utils.ErrorResponse(c, 500, "Database error")
		return
	}

	utils.SuccessResponse(c, 200, "Business profile retrieved successfully", businessProfile)
}

// DeleteBusinessProfile deletes the business profile for the authenticated user
func DeleteBusinessProfile(c *gin.Context) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	businessCollection := database.DB.Collection("business_profiles")

	result, err := businessCollection.DeleteOne(ctx, bson.M{"user_id": userObjID})
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete business profile")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, 404, "Business profile not found")
		return
	}

	utils.SuccessResponse(c, 200, "Business profile deleted successfully", nil)
}

