package main

import (
	"log"
	"sheleads-backend/config"
	"sheleads-backend/database"
	"sheleads-backend/handlers"
	"sheleads-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB
	database.Connect()

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "SheLeads Backend",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
		}

		// User profile
		v1.GET("/profile", handlers.GetProfile)

		// Work logging
		work := v1.Group("/work")
		{
			work.POST("/log", handlers.LogWork)
			work.GET("/logs", handlers.GetWorkLogs)
			work.GET("/logs/:id", handlers.GetWorkLogByID)
			work.DELETE("/logs/:id", handlers.DeleteWorkLog)
		}

		// Analytics
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/summary", handlers.GetAnalyticsSummary)
			analytics.GET("/stats", handlers.GetActivityStats)
			analytics.GET("/burnout", handlers.GetBurnoutStatus)
			analytics.GET("/calendar-insights", handlers.GetCalendarInsights)
			analytics.GET("/market-value", handlers.GetMarketValue)
		}

		// Reports
		reports := v1.Group("/reports")
		{
			reports.GET("/monthly/pdf", handlers.GenerateMonthlyReport)
			reports.GET("/monthly/data", handlers.GetMonthlyReportData)
		}

		// Business Profile (Marketing Automation)
		business := v1.Group("/business")
		{
			business.POST("/profile", handlers.CreateBusinessProfile)
			business.GET("/profile", handlers.GetBusinessProfile)
			business.DELETE("/profile", handlers.DeleteBusinessProfile)
		}

		// Content Generation (Marketing Automation)
		content := v1.Group("/content")
		{
			content.POST("/blog", handlers.GenerateBlog)
			content.POST("/social", handlers.GenerateSocialPost)
			content.GET("", handlers.GetAllContent)
			content.GET("/:id", handlers.GetContentByID)
			content.DELETE("/:id", handlers.DeleteContent)
			content.PUT("/:id/status", handlers.UpdateContentStatus)
			content.POST("/:id/generate-image", handlers.GenerateImageForContent)
			content.POST("/:id/post", handlers.PostContentToSocial)
		}

		// Content Metrics (Marketing Automation)
		metrics := v1.Group("/metrics")
		{
			metrics.POST("", handlers.AddMetrics)
			metrics.GET("/content/:content_id", handlers.GetMetricsByContentID)
			metrics.GET("/aggregated", handlers.GetAggregatedMetrics)
		}

		// Asset Management (Marketing Automation)
		assets := v1.Group("/assets")
		{
			assets.POST("/upload", handlers.UploadAsset)
			assets.GET("", handlers.GetUserAssets)
		}
	}

	// Start server
	port := config.AppConfig.Port
	log.Printf("🚀 SheLeads Backend starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

