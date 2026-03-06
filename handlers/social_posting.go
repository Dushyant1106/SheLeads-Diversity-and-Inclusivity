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
	"go.mongodb.org/mongo-driver/mongo"
)

// GenerateImageForContent generates an image using Runway AI for a piece of content
func GenerateImageForContent(c *gin.Context) {
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
		ImagePrompt      string   `json:"image_prompt" binding:"required"`
		ReferenceImages  []string `json:"reference_images"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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

	// Get business profile to add context to image generation
	businessCollection := database.DB.Collection("business_profiles")
	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"_id": content.BusinessID}).Decode(&businessProfile)
	if err != nil && err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 500, "Failed to fetch business profile")
		return
	}

	// Enhance the image prompt with business context
	enhancedPrompt := input.ImagePrompt
	if businessProfile.BusinessName != "" {
		enhancedPrompt = fmt.Sprintf("%s. Business context: %s in %s industry, targeting %s. Style: professional, high-quality, relevant to %s",
			input.ImagePrompt,
			businessProfile.BusinessName,
			businessProfile.Industry,
			businessProfile.TargetAudience,
			content.Title,
		)
	}

	// Fetch user's uploaded assets to use as reference images
	referenceImages := input.ReferenceImages
	if businessProfile.Metadata != nil {
		if assets, ok := businessProfile.Metadata["assets"].(map[string]interface{}); ok {
			// Get reference images from assets
			if refImages, ok := assets["reference_image"].([]interface{}); ok {
				for _, asset := range refImages {
					if assetMap, ok := asset.(map[string]interface{}); ok {
						if url, ok := assetMap["url"].(string); ok {
							referenceImages = append(referenceImages, url)
						}
					}
				}
			}
			// Also include logo if available
			if businessProfile.LogoURL != "" {
				referenceImages = append(referenceImages, businessProfile.LogoURL)
			}
		}
	}

	// Limit to 5 reference images (Runway API limit)
	if len(referenceImages) > 5 {
		referenceImages = referenceImages[:5]
	}

	// Generate image using Runway
	runwayService, err := services.NewRunwayService()
	if err != nil {
		utils.ErrorResponse(c, 500, "Runway service unavailable: "+err.Error())
		return
	}

	imagePath, err := runwayService.GenerateImage(ctx, enhancedPrompt, referenceImages)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate image: "+err.Error())
		return
	}

	// Upload generated image to S3
	s3Service, err := services.NewS3Service()
	if err != nil {
		utils.ErrorResponse(c, 500, "S3 service unavailable: "+err.Error())
		return
	}

	s3Key := "generated-images/" + contentID + ".jpg"
	s3URL, err := s3Service.UploadFileFromPath(ctx, imagePath, s3Key, "image/jpeg")
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to upload image to S3: "+err.Error())
		return
	}

	// Update content with image URL
	_, err = contentCollection.UpdateOne(ctx,
		bson.M{"_id": contentObjID},
		bson.M{
			"$set": bson.M{
				"image_url":  s3URL,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update content")
		return
	}

	utils.SuccessResponse(c, 200, "Image generated successfully", gin.H{
		"image_url":   s3URL,
		"local_path":  imagePath,
		"content_id":  contentID,
	})
}

// PostContentToSocial posts content to social media platforms
func PostContentToSocial(c *gin.Context) {
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
		Platforms []string `json:"platforms" binding:"required"`
		ImagePath string   `json:"image_path"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Invalid input: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Get content
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

	// Initialize social media service
	socialService, err := services.NewSocialMediaService()
	if err != nil {
		utils.ErrorResponse(c, 500, "Social media service unavailable: "+err.Error())
		return
	}

	// Post to each platform
	postURLs := make(map[string]string)
	errors := make(map[string]string)

	for _, platform := range input.Platforms {
		postURL, err := socialService.PostToSocial(platform, content.Content, input.ImagePath)
		if err != nil {
			errors[platform] = err.Error()
		} else {
			postURLs[platform] = postURL
		}
	}

	// Update content status
	if len(postURLs) > 0 {
		_, err = contentCollection.UpdateOne(ctx,
			bson.M{"_id": contentObjID},
			bson.M{
				"$set": bson.M{
					"status":     "posted",
					"post_urls":  postURLs,
					"posted_at":  time.Now(),
					"updated_at": time.Now(),
				},
			},
		)
	}

	if len(errors) > 0 && len(postURLs) == 0 {
		utils.ErrorResponse(c, 500, "Failed to post to any platform")
		return
	}

	utils.SuccessResponse(c, 200, "Content posted to social media", gin.H{
		"post_urls": postURLs,
		"errors":    errors,
	})
}

