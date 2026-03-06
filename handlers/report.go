package handlers

import (
	"sheleads-backend/middleware"
	"sheleads-backend/services"
	"sheleads-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateMonthlyReport(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		utils.ErrorResponse(c, 401, "Unauthorized")
		return
	}

	// Get month and year from query parameters
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		utils.ErrorResponse(c, 400, "Invalid month. Must be between 1 and 12")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		utils.ErrorResponse(c, 400, "Invalid year")
		return
	}

	// Generate PDF report
	pdfService := services.NewPDFService()
	filename, err := pdfService.GenerateMonthlyReport(userIDStr, time.Month(month), year)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate report: "+err.Error())
		return
	}

	// Return the file
	c.File(filename)
}

func GetMonthlyReportData(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, 400, "User ID is required")
		return
	}

	// Get month and year from query parameters
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		utils.ErrorResponse(c, 400, "Invalid month. Must be between 1 and 12")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid year")
		return
	}

	// Calculate date range
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Use the existing activity stats function
	c.Set("start_date", startDate.Format("2006-01-02"))
	c.Set("end_date", endDate.Format("2006-01-02"))

	GetActivityStats(c)
}

