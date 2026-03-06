package services

import (
	"context"
	"sheleads-backend/config"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BurnoutService struct {
	twilioService *TwilioService
}

func NewBurnoutService() *BurnoutService {
	return &BurnoutService{
		twilioService: NewTwilioService(),
	}
}

type BurnoutStatus struct {
	IsBurnout       bool    `json:"is_burnout"`
	HoursInWindow   float64 `json:"hours_in_window"`
	DaysAnalyzed    int     `json:"days_analyzed"`
	Threshold       float64 `json:"threshold"`
	AlertSent       bool    `json:"alert_sent"`
	Recommendation  string  `json:"recommendation"`
}

func (b *BurnoutService) CheckBurnout(userIDStr string) (*BurnoutStatus, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user
	usersCollection := database.DB.Collection("users")
	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -config.AppConfig.BurnoutDaysWindow)

	// Get work logs in the window
	worklogsCollection := database.DB.Collection("worklogs")
	cursor, err := worklogsCollection.Find(ctx, bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var workLogs []models.WorkLog
	if err = cursor.All(ctx, &workLogs); err != nil {
		return nil, err
	}

	// Calculate total hours
	var totalHours float64
	for _, log := range workLogs {
		totalHours += log.EstimatedHours
	}

	// Check if burnout threshold exceeded
	threshold := config.AppConfig.BurnoutHoursThreshold * float64(config.AppConfig.BurnoutDaysWindow)
	isBurnout := totalHours > threshold

	status := &BurnoutStatus{
		IsBurnout:     isBurnout,
		HoursInWindow: totalHours,
		DaysAnalyzed:  config.AppConfig.BurnoutDaysWindow,
		Threshold:     threshold,
		AlertSent:     false,
	}

	// Generate recommendation
	if isBurnout {
		status.Recommendation = "You're working too hard! Please take time to rest and care for yourself."

		// Send alert to emergency contact
		err := b.twilioService.SendBurnoutAlert(user.EmergencyContact, user.Name, totalHours)
		if err == nil {
			status.AlertSent = true
		}
	} else {
		percentOfThreshold := (totalHours / threshold) * 100
		if percentOfThreshold > 75 {
			status.Recommendation = "You're approaching burnout levels. Please consider taking breaks."
		} else {
			status.Recommendation = "Your work levels are healthy. Keep taking care of yourself!"
		}
	}

	return status, nil
}

