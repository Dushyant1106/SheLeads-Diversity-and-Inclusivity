package handlers

import (
	"context"
	"path/filepath"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"sheleads-backend/services"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadAsset uploads an asset (image, logo, etc.) to S3
func UploadAsset(c *gin.Context) {
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

	// Get asset type from query parameter
	assetType := c.DefaultQuery("asset_type", "general")
	
	// Validate asset type
	validTypes := map[string]bool{
		"logo":            true,
		"item":            true,
		"reference_image": true,
		"general":         true,
	}
	if !validTypes[assetType] {
		utils.ErrorResponse(c, 400, "Invalid asset_type. Allowed: logo, item, reference_image, general")
		return
	}

	// Get uploaded file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, "No file uploaded: "+err.Error())
		return
	}
	defer file.Close()

	// Validate file extension
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}
	ext := filepath.Ext(fileHeader.Filename)
	if !allowedExts[ext] {
		utils.ErrorResponse(c, 400, "File type not allowed. Allowed: jpg, jpeg, png, gif, webp, svg")
		return
	}

	// Check file size (10MB limit)
	maxSize := int64(10 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		utils.ErrorResponse(c, 400, "File too large. Maximum size: 10MB")
		return
	}

	// Upload to S3
	s3Service, err := services.NewS3Service()
	if err != nil {
		utils.ErrorResponse(c, 500, "S3 service unavailable: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s3URL, err := s3Service.UploadFile(ctx, file, fileHeader, userIDStr, assetType)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to upload file: "+err.Error())
		return
	}

	// Update business profile if it's a logo
	if assetType == "logo" {
		businessCollection := database.DB.Collection("business_profiles")
		_, err = businessCollection.UpdateOne(ctx,
			bson.M{"user_id": userObjID},
			bson.M{
				"$set": bson.M{
					"logo_url":   s3URL,
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			// Log error but don't fail the upload
			utils.ErrorResponse(c, 500, "File uploaded but failed to update business profile")
			return
		}
	}

	// Store asset metadata in business profile
	businessCollection := database.DB.Collection("business_profiles")
	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&businessProfile)
	
	if err == nil {
		// Add asset to metadata
		if businessProfile.Metadata == nil {
			businessProfile.Metadata = make(map[string]interface{})
		}
		
		assets, ok := businessProfile.Metadata["assets"].(map[string]interface{})
		if !ok {
			assets = make(map[string]interface{})
		}
		
		assetList, ok := assets[assetType].([]interface{})
		if !ok {
			assetList = []interface{}{}
		}
		
		assetInfo := map[string]interface{}{
			"url":         s3URL,
			"filename":    fileHeader.Filename,
			"uploaded_at": time.Now(),
			"size":        fileHeader.Size,
		}
		
		assetList = append(assetList, assetInfo)
		
		// Keep only last 20 assets per type
		if len(assetList) > 20 {
			assetList = assetList[len(assetList)-20:]
		}
		
		assets[assetType] = assetList
		businessProfile.Metadata["assets"] = assets
		
		_, err = businessCollection.UpdateOne(ctx,
			bson.M{"user_id": userObjID},
			bson.M{
				"$set": bson.M{
					"metadata":   businessProfile.Metadata,
					"updated_at": time.Now(),
				},
			},
		)
	}

	utils.SuccessResponse(c, 201, "File uploaded successfully", gin.H{
		"s3_url":     s3URL,
		"asset_type": assetType,
		"filename":   fileHeader.Filename,
		"size":       fileHeader.Size,
	})
}

// GetUserAssets retrieves user's uploaded assets
func GetUserAssets(c *gin.Context) {
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

	assetType := c.Query("asset_type")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	businessCollection := database.DB.Collection("business_profiles")
	var businessProfile models.BusinessProfile
	err = businessCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&businessProfile)
	
	if err != nil {
		utils.SuccessResponse(c, 200, "No assets found", gin.H{
			"assets":   []interface{}{},
			"logo_url": nil,
		})
		return
	}

	logoURL := businessProfile.LogoURL
	
	if businessProfile.Metadata == nil {
		utils.SuccessResponse(c, 200, "Assets retrieved successfully", gin.H{
			"assets":   []interface{}{},
			"logo_url": logoURL,
		})
		return
	}

	assets, ok := businessProfile.Metadata["assets"].(map[string]interface{})
	if !ok {
		assets = make(map[string]interface{})
	}

	// Filter by asset_type if specified
	if assetType != "" {
		filteredAssets := assets[assetType]
		if filteredAssets == nil {
			filteredAssets = []interface{}{}
		}
		utils.SuccessResponse(c, 200, "Assets retrieved successfully", gin.H{
			"assets":   filteredAssets,
			"logo_url": logoURL,
		})
		return
	}

	// Return all assets
	utils.SuccessResponse(c, 200, "Assets retrieved successfully", gin.H{
		"assets":   assets,
		"logo_url": logoURL,
	})
}

