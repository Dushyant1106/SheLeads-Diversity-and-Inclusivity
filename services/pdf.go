package services

import (
	"context"
	"fmt"
	"sheleads-backend/database"
	"sheleads-backend/models"
	"time"

	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

type ReportData struct {
	User              models.User
	WorkLogs          []models.WorkLog
	TotalHours        float64
	TotalPoints       int
	StartDate         time.Time
	EndDate           time.Time
	CategoryBreakdown map[models.WorkCategory]CategoryStats
	TotalMarketValue  float64
	BusiestDay        string
	AvgHoursPerDay    float64
	PeakTime          string
}

type CategoryStats struct {
	Count       int
	Hours       float64
	Points      int
	MarketValue float64
}

func (p *PDFService) GenerateMonthlyReport(userIDStr string, month time.Month, year int) (string, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user
	usersCollection := database.DB.Collection("users")
	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return "", err
	}

	// Calculate date range for the month
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Get work logs for the month
	worklogsCollection := database.DB.Collection("worklogs")
	cursor, err := worklogsCollection.Find(ctx, bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	var workLogs []models.WorkLog
	if err = cursor.All(ctx, &workLogs); err != nil {
		return "", err
	}

	// Calculate statistics
	reportData := p.calculateStats(user, workLogs, startDate, endDate)

	// Generate PDF
	return p.createPDF(reportData)
}

func (p *PDFService) calculateStats(user models.User, workLogs []models.WorkLog, startDate, endDate time.Time) ReportData {
	var totalHours float64
	var totalPoints int
	var totalMarketValue float64
	categoryBreakdown := make(map[models.WorkCategory]CategoryStats)
	dayHours := make(map[string]float64) // Track hours per day of week

	// Indian market rates (₹ per hour)
	marketRates := map[models.WorkCategory]float64{
		models.Cooking:         300,
		models.Cleaning:        250,
		models.Childcare:       200,
		models.Eldercare:       350,
		models.Laundry:         150,
		models.Shopping:        100,
		models.Gardening:       200,
		models.HomeMaintenance: 400,
		models.Other:           200,
	}

	for _, log := range workLogs {
		totalHours += log.EstimatedHours
		totalPoints += log.Points

		// Calculate market value
		rate := marketRates[log.Category]
		value := log.EstimatedHours * rate
		totalMarketValue += value

		// Track day of week
		dayOfWeek := log.CreatedAt.Weekday().String()
		dayHours[dayOfWeek] += log.EstimatedHours

		stats := categoryBreakdown[log.Category]
		stats.Count++
		stats.Hours += log.EstimatedHours
		stats.Points += log.Points
		stats.MarketValue += value
		categoryBreakdown[log.Category] = stats
	}

	// Find busiest day
	busiestDay := "N/A"
	maxHours := 0.0
	for day, hours := range dayHours {
		if hours > maxHours {
			maxHours = hours
			busiestDay = day
		}
	}

	// Calculate average hours per day
	avgHoursPerDay := 0.0
	if len(workLogs) > 0 {
		daysInPeriod := endDate.Sub(startDate).Hours() / 24
		if daysInPeriod > 0 {
			avgHoursPerDay = totalHours / daysInPeriod
		}
	}

	// Determine peak time (simplified)
	peakTime := "Daytime"
	morningHours := 0.0
	afternoonHours := 0.0
	eveningHours := 0.0

	for _, log := range workLogs {
		hour := log.CreatedAt.Hour()
		if hour >= 6 && hour < 12 {
			morningHours += log.EstimatedHours
		} else if hour >= 12 && hour < 18 {
			afternoonHours += log.EstimatedHours
		} else {
			eveningHours += log.EstimatedHours
		}
	}

	if morningHours > afternoonHours && morningHours > eveningHours {
		peakTime = "Morning"
	} else if eveningHours > morningHours && eveningHours > afternoonHours {
		peakTime = "Evening"
	}

	return ReportData{
		User:              user,
		WorkLogs:          workLogs,
		TotalHours:        totalHours,
		TotalPoints:       totalPoints,
		StartDate:         startDate,
		EndDate:           endDate,
		CategoryBreakdown: categoryBreakdown,
		TotalMarketValue:  totalMarketValue,
		BusiestDay:        busiestDay,
		AvgHoursPerDay:    avgHoursPerDay,
		PeakTime:          peakTime,
	}
}

func (p *PDFService) createPDF(data ReportData) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(139, 0, 139) // Purple color
	pdf.CellFormat(190, 15, "SheLeads Care Work Report", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// User Info Section
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 10, "Personal Information", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("Name: %s", data.User.Name), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Email: %s", data.User.Email), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Report Period: %s to %s", 
		data.StartDate.Format("Jan 02, 2006"), 
		data.EndDate.Format("Jan 02, 2006")), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Summary Section
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("Total Work Entries: %d", len(data.WorkLogs)), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Total Hours Worked: %.1f hours", data.TotalHours), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Total Points Earned: %d points", data.TotalPoints), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 128, 0) // Green for market value
	pdf.CellFormat(190, 8, fmt.Sprintf("Total Market Value: Rs %.2f", data.TotalMarketValue), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0) // Reset to black
	pdf.Ln(5)

	// Work Pattern Insights Section
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "Work Pattern Insights", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("Busiest Day: %s", data.BusiestDay), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Average Hours Per Day: %.1f hours", data.AvgHoursPerDay), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Peak Work Time: %s", data.PeakTime), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Category Breakdown with Market Value
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "Work Category Breakdown", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 8, "Category", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Count", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Hours", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Points", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Value (Rs)", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	for category, stats := range data.CategoryBreakdown {
		pdf.CellFormat(50, 8, string(category), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", stats.Count), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.1f", stats.Hours), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%d", stats.Points), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.0f", stats.MarketValue), "1", 1, "C", false, 0, "")
	}
	pdf.Ln(5)

	// Loan Eligibility Section
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "Loan Eligibility", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 12)

	// Calculate discount based on points
	discountPercent := float64(data.TotalPoints) / 1000.0 * 10.0
	if discountPercent > 50 {
		discountPercent = 50
	}

	loanTiers := []struct {
		Name      string
		Amount    int
		BaseRate  float64
		MinPoints int
	}{
		{"Micro Business Starter", 50000, 14.0, 100},
		{"Small Business Growth", 200000, 12.0, 500},
		{"Business Expansion", 500000, 11.0, 1000},
		{"Enterprise Growth", 1000000, 10.0, 2000},
	}

	for _, tier := range loanTiers {
		if data.TotalPoints >= tier.MinPoints {
			reducedRate := tier.BaseRate * (1 - discountPercent/100)
			pdf.SetTextColor(0, 128, 0) // Green for eligible
			pdf.CellFormat(190, 8, fmt.Sprintf("✓ %s: Rs %d @ %.2f%% (was %.1f%%)",
				tier.Name, tier.Amount, reducedRate, tier.BaseRate), "", 1, "L", false, 0, "")
		} else {
			pdf.SetTextColor(128, 128, 128) // Gray for not eligible
			pdf.CellFormat(190, 8, fmt.Sprintf("  %s: Rs %d (Need %d points)",
				tier.Name, tier.Amount, tier.MinPoints), "", 1, "L", false, 0, "")
		}
	}
	pdf.SetTextColor(0, 0, 0) // Reset to black
	pdf.Ln(5)

	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(190, 10, "Your unpaid care work is valuable and recognized.", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 10, "Generated by SheLeads Platform", "", 1, "C", false, 0, "")

	// Save PDF
	filename := fmt.Sprintf("uploads/report_%d_%s.pdf", data.User.ID, time.Now().Format("20060102_150405"))
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

