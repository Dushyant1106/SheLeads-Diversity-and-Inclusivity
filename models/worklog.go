package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkCategory string

const (
	CategoryCooking         WorkCategory = "cooking"
	CategoryCleaning        WorkCategory = "cleaning"
	CategoryChildcare       WorkCategory = "childcare"
	CategoryEldercare       WorkCategory = "eldercare"
	CategoryLaundry         WorkCategory = "laundry"
	CategoryShopping        WorkCategory = "shopping"
	CategoryGardening       WorkCategory = "gardening"
	CategoryHomeMaintenance WorkCategory = "home_maintenance"
	CategoryOther           WorkCategory = "other"
)

type WorkLog struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	Category        WorkCategory       `json:"category" bson:"category"`
	Description     string             `json:"description" bson:"description"`
	ImagePath       string             `json:"image_path" bson:"image_path"`
	ConfidenceScore float64            `json:"confidence_score" bson:"confidence_score"`
	AIVerification  string             `json:"ai_verification" bson:"ai_verification"`
	Points          int                `json:"points" bson:"points"`
	EstimatedHours  float64            `json:"estimated_hours" bson:"estimated_hours"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// CalculatePoints calculates points based on confidence score and category
func (w *WorkLog) CalculatePoints() {
	basePoints := 10
	
	// Category multipliers
	categoryMultiplier := map[WorkCategory]float64{
		CategoryChildcare:       1.5,
		CategoryEldercare:       1.5,
		CategoryCooking:         1.2,
		CategoryCleaning:        1.0,
		CategoryLaundry:         1.0,
		CategoryShopping:        0.8,
		CategoryGardening:       1.1,
		CategoryHomeMaintenance: 1.3,
		CategoryOther:           1.0,
	}
	
	multiplier := categoryMultiplier[w.Category]
	if multiplier == 0 {
		multiplier = 1.0
	}
	
	// Points = basePoints * categoryMultiplier * confidenceScore * estimatedHours
	w.Points = int(float64(basePoints) * multiplier * w.ConfidenceScore * w.EstimatedHours)
	
	// Mark as verified if confidence score is high enough
	if w.ConfidenceScore >= 0.7 {
		w.Verified = true
	}
}

