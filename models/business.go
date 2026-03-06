package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BusinessProfile represents a woman's business idea/profile
type BusinessProfile struct {
	ID                   primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	UserID               primitive.ObjectID     `bson:"user_id" json:"user_id"`
	BusinessName         string                 `bson:"business_name" json:"business_name" binding:"required"`
	Industry             string                 `bson:"industry" json:"industry" binding:"required"`
	Location             string                 `bson:"location" json:"location"`
	Description          string                 `bson:"description" json:"description"`
	TargetAudience       string                 `bson:"target_audience" json:"target_audience"`
	UniqueSellingPoints  []string               `bson:"unique_selling_points" json:"unique_selling_points"`
	Website              string                 `bson:"website" json:"website"`
	LogoURL              string                 `bson:"logo_url" json:"logo_url"`
	SocialMediaHandles   map[string]string      `bson:"social_media_handles" json:"social_media_handles"` // platform -> handle
	Metadata             map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`      // For storing assets and other metadata
	CreatedAt            time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time              `bson:"updated_at" json:"updated_at"`
}

// BusinessProfileInput for creating/updating business profile
type BusinessProfileInput struct {
	BusinessName        string            `json:"business_name" binding:"required"`
	Industry            string            `json:"industry" binding:"required"`
	Location            string            `json:"location"`
	Description         string            `json:"description"`
	TargetAudience      string            `json:"target_audience"`
	UniqueSellingPoints []string          `json:"unique_selling_points"`
	Website             string            `json:"website"`
	SocialMediaHandles  map[string]string `json:"social_media_handles"`
}

