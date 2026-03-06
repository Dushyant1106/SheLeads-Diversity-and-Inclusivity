package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeBlog   ContentType = "blog"
	ContentTypeSocial ContentType = "social"
)

// ContentStatus represents the status of content
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusGenerated ContentStatus = "generated"
	ContentStatusPosted    ContentStatus = "posted"
	ContentStatusFailed    ContentStatus = "failed"
)

// GeneratedContent represents a blog post or social media post
type GeneratedContent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	BusinessID  primitive.ObjectID `bson:"business_id" json:"business_id"`
	Type        ContentType        `bson:"type" json:"type"`
	Status      ContentStatus      `bson:"status" json:"status"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Platform    string             `bson:"platform,omitempty" json:"platform,omitempty"` // For social: linkedin, twitter, instagram, facebook
	Hashtags    []string           `bson:"hashtags,omitempty" json:"hashtags,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	PostURL     string             `bson:"post_url,omitempty" json:"post_url,omitempty"` // URL after posting
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	PostedAt    *time.Time         `bson:"posted_at,omitempty" json:"posted_at,omitempty"`
}

// BlogGenerationRequest for generating blog posts
type BlogGenerationRequest struct {
	Topic       string   `json:"topic" binding:"required"`
	Keywords    []string `json:"keywords"`
	Tone        string   `json:"tone"` // professional, casual, friendly, etc.
	Length      string   `json:"length"` // short, medium, long
}

// SocialPostGenerationRequest for generating social media posts
type SocialPostGenerationRequest struct {
	Topic       string   `json:"topic" binding:"required"`
	Platforms   []string `json:"platforms" binding:"required"` // linkedin, twitter, instagram, facebook
	Tone        string   `json:"tone"`
	Hashtags    []string `json:"hashtags"`
	ImageStyle  string   `json:"image_style"` // For future image generation
}

// ContentMetrics represents engagement metrics for content
type ContentMetrics struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ContentID     primitive.ObjectID `bson:"content_id" json:"content_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Platform      string             `bson:"platform" json:"platform"`
	Likes         int                `bson:"likes" json:"likes"`
	Comments      int                `bson:"comments" json:"comments"`
	Shares        int                `bson:"shares" json:"shares"`
	Impressions   int                `bson:"impressions" json:"impressions"`
	Reach         int                `bson:"reach" json:"reach"`
	EngagementRate float64           `bson:"engagement_rate" json:"engagement_rate"`
	CollectedAt   time.Time          `bson:"collected_at" json:"collected_at"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

// MetricsInput for manually adding/updating metrics
type MetricsInput struct {
	ContentID   string `json:"content_id" binding:"required"`
	Likes       int    `json:"likes"`
	Comments    int    `json:"comments"`
	Shares      int    `json:"shares"`
	Impressions int    `json:"impressions"`
	Reach       int    `json:"reach"`
}

