package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                      string
	MongoURL                  string
	DBName                    string
	JWTSecret                 string
	GeminiAPIKey              string
	TwilioAccountSID          string
	TwilioAuthToken           string
	TwilioMessagingServiceSID string
	BurnoutHoursThreshold     float64
	BurnoutDaysWindow         int
	// AWS S3 Configuration
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSS3BucketName    string
	AWSRegion          string
	// Runway AI Configuration
	RunwayAPIKey string
	// Social Media Configuration
	TwitterAPIKey             string
	TwitterAPISecret          string
	TwitterAccessToken        string
	TwitterAccessTokenSecret  string
	InstagramUsername         string
	InstagramPassword         string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	burnoutHours, _ := strconv.ParseFloat(getEnv("BURNOUT_HOURS_THRESHOLD", "12"), 64)
	burnoutDays, _ := strconv.Atoi(getEnv("BURNOUT_DAYS_WINDOW", "7"))

	AppConfig = &Config{
		Port:                      getEnv("PORT", "8080"),
		MongoURL:                  getEnv("MONGO_URL", ""),
		DBName:                    getEnv("DB_NAME", "caretrack"),
		JWTSecret:                 getEnv("JWT_SECRET", ""),
		GeminiAPIKey:              getEnv("GEMINI_API_KEY", ""),
		TwilioAccountSID:          getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:           getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioMessagingServiceSID: getEnv("TWILIO_MESSAGING_SERVICE_SID", ""),
		BurnoutHoursThreshold:     burnoutHours,
		BurnoutDaysWindow:         burnoutDays,
		// AWS S3
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3BucketName:    getEnv("AWS_S3_BUCKET_NAME", ""),
		AWSRegion:          getEnv("AWS_REGION", "us-east-1"),
		// Runway AI
		RunwayAPIKey: getEnv("RUNWAY_API_KEY", ""),
		// Social Media
		TwitterAPIKey:            getEnv("TWITTER_API_KEY", ""),
		TwitterAPISecret:         getEnv("TWITTER_API_SECRET", ""),
		TwitterAccessToken:       getEnv("TWITTER_ACCESS_TOKEN", ""),
		TwitterAccessTokenSecret: getEnv("TWITTER_ACCESS_TOKEN_SECRET", ""),
		InstagramUsername:        getEnv("INSTAGRAM_USERNAME", ""),
		InstagramPassword:        getEnv("INSTAGRAM_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

