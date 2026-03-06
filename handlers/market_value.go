package handlers

import (
	"context"
	"math"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"sheleads-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Market rates per hour for different categories (in INR - Indian Rupees)
// Based on Urban Company, UrbanClap, and other Indian domestic service platforms
var marketRates = map[models.WorkCategory]float64{
	models.CategoryCooking:         300.0,  // Cook/Chef rates (₹250-350/hr)
	models.CategoryCleaning:        250.0,  // Professional cleaning service (₹200-300/hr)
	models.CategoryChildcare:       200.0,  // Nanny/babysitter rates (₹150-250/hr)
	models.CategoryEldercare:       350.0,  // Elder care professional (₹300-400/hr)
	models.CategoryLaundry:         150.0,  // Laundry/ironing service (₹100-200/hr)
	models.CategoryShopping:        100.0,  // Personal shopper/errand service (₹80-120/hr)
	models.CategoryGardening:       200.0,  // Gardening/landscaping service (₹150-250/hr)
	models.CategoryHomeMaintenance: 400.0,  // Handyman/repair service (₹350-450/hr)
	models.CategoryOther:           200.0,  // General household help (₹150-250/hr)
}

type CategoryBreakdown struct {
	Category string  `json:"category"`
	Hours    float64 `json:"hours"`
	Rate     float64 `json:"rate"`
	Value    float64 `json:"value"`
}

type LoanOption struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Term            string  `json:"term"`
	OriginalRate    float64 `json:"original_rate"`
	ReducedRate     float64 `json:"reduced_rate"`
	DiscountPercent float64 `json:"discount_percent"`
	MonthlyPayment  float64 `json:"monthly_payment"`
	TotalSavings    float64 `json:"total_savings"`
	Eligible        bool    `json:"eligible"`
}

type MarketValueResponse struct {
	TotalMarketValue  float64             `json:"total_market_value"`
	CategoryBreakdown []CategoryBreakdown `json:"category_breakdown"`
	TotalPoints       int                 `json:"total_points"`
	LoanOptions       []LoanOption        `json:"loan_options"`
}

// GetMarketValue calculates the market value of logged work and loan options
func GetMarketValue(c *gin.Context) {
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

	// Get all work logs
	worklogsCollection := database.DB.Collection("worklogs")
	cursor, err := worklogsCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch work logs")
		return
	}
	defer cursor.Close(ctx)

	var workLogs []models.WorkLog
	if err = cursor.All(ctx, &workLogs); err != nil {
		utils.ErrorResponse(c, 500, "Failed to parse work logs")
		return
	}

	// Calculate market value by category
	categoryTotals := make(map[models.WorkCategory]float64)
	totalPoints := 0

	for _, log := range workLogs {
		categoryTotals[log.Category] += log.EstimatedHours
		totalPoints += log.Points
	}

	// Build category breakdown
	var breakdown []CategoryBreakdown
	totalMarketValue := 0.0

	for category, hours := range categoryTotals {
		rate := marketRates[category]
		value := hours * rate
		totalMarketValue += value

		breakdown = append(breakdown, CategoryBreakdown{
			Category: string(category),
			Hours:    hours,
			Rate:     rate,
			Value:    value,
		})
	}

	// Calculate loan options with point-based discounts
	loanOptions := calculateLoanOptions(totalPoints, totalMarketValue)

	response := MarketValueResponse{
		TotalMarketValue:  totalMarketValue,
		CategoryBreakdown: breakdown,
		TotalPoints:       totalPoints,
		LoanOptions:       loanOptions,
	}

	utils.SuccessResponse(c, 200, "Market value calculated", response)
}

func calculateLoanOptions(points int, marketValue float64) []LoanOption {
	// Calculate discount based on points (max 50% discount)
	discountPercent := math.Min(float64(points)/1000.0*10.0, 50.0)

	// Loan amounts in Indian Rupees (INR)
	loans := []LoanOption{
		{
			Name:         "Micro Business Starter Loan",
			Description:  "Perfect for launching your small business or home-based venture",
			Amount:       50000,  // ₹50,000
			Term:         "12 months",
			OriginalRate: 14.0,  // Typical Indian microfinance rate
			Eligible:     points >= 100,
		},
		{
			Name:         "Small Business Growth Loan",
			Description:  "Expand your existing business operations and inventory",
			Amount:       200000,  // ₹2,00,000
			Term:         "24 months",
			OriginalRate: 12.0,  // Small business loan rate
			Eligible:     points >= 500,
		},
		{
			Name:         "Business Expansion Loan",
			Description:  "Scale your business with equipment, staff, and marketing",
			Amount:       500000,  // ₹5,00,000
			Term:         "36 months",
			OriginalRate: 11.0,  // Business expansion rate
			Eligible:     points >= 1000,
		},
		{
			Name:         "Enterprise Growth Loan",
			Description:  "Major expansion for established businesses - new location, franchise, etc.",
			Amount:       1000000,  // ₹10,00,000
			Term:         "48 months",
			OriginalRate: 10.0,  // Enterprise loan rate
			Eligible:     points >= 2000,
		},
	}

	// Apply discount and calculate payments
	for i := range loans {
		loans[i].ReducedRate = loans[i].OriginalRate * (1 - discountPercent/100)
		loans[i].DiscountPercent = discountPercent

		// Calculate monthly payment using loan formula
		monthlyRate := loans[i].ReducedRate / 100 / 12
		var months float64
		switch loans[i].Term {
		case "12 months":
			months = 12
		case "24 months":
			months = 24
		case "36 months":
			months = 36
		case "48 months":
			months = 48
		}

		if monthlyRate > 0 {
			loans[i].MonthlyPayment = loans[i].Amount * monthlyRate * math.Pow(1+monthlyRate, months) / (math.Pow(1+monthlyRate, months) - 1)
		} else {
			loans[i].MonthlyPayment = loans[i].Amount / months
		}

		// Calculate total savings
		originalMonthlyRate := loans[i].OriginalRate / 100 / 12
		var originalPayment float64
		if originalMonthlyRate > 0 {
			originalPayment = loans[i].Amount * originalMonthlyRate * math.Pow(1+originalMonthlyRate, months) / (math.Pow(1+originalMonthlyRate, months) - 1)
		} else {
			originalPayment = loans[i].Amount / months
		}
		loans[i].TotalSavings = (originalPayment - loans[i].MonthlyPayment) * months
	}

	return loans
}

