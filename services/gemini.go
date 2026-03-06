package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sheleads-backend/config"
	"sheleads-backend/models"
	"strings"
	

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client *genai.Client
}

type VerificationResult struct {
	IsValid         bool    `json:"is_valid"`
	ConfidenceScore float64 `json:"confidence_score"`
	Explanation     string  `json:"explanation"`
	EstimatedHours  float64 `json:"estimated_hours"`
}

func NewGeminiService(ctx context.Context) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.AppConfig.GeminiAPIKey))
	if err != nil {
		return nil, err
	}
	return &GeminiService{client: client}, nil
}

func (g *GeminiService) VerifyWork(ctx context.Context, imagePath string, category models.WorkCategory, description string) (*VerificationResult, error) {
	// Read the image file
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	// Create the prompt for verification
	prompt := g.createVerificationPrompt(category, description)

	// Get the generative model - using gemini-2.0-flash-exp which is faster
	model := g.client.GenerativeModel("gemini-2.5-flash")

	// Set generation config for faster responses
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(1024) // Increased to ensure full JSON response

	// Prepare parts for Gemini - using the correct API
	parts := []genai.Part{
		genai.ImageData("jpeg", imageBytes),
		genai.Text(prompt),
	}

	// Generate content using Gemini with context timeout
	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini AI")
	}

	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Parse the response
	return g.parseVerificationResponse(responseText)
}

func (g *GeminiService) createVerificationPrompt(category models.WorkCategory, description string) string {
	return fmt.Sprintf(`You are an AI assistant helping to verify unpaid care work for a women's empowerment platform.

Task Category: %s
User Description: %s

Please analyze the provided image and determine if it genuinely represents the claimed unpaid care work.

Respond ONLY with a valid JSON object in this exact format:
{
  "is_valid": true/false,
  "confidence_score": 0.0-1.0,
  "explanation": "brief explanation of your assessment",
  "estimated_hours": 0.5-8.0
}

Criteria for verification:
1. Does the image show evidence of the claimed work category?
2. Is the image authentic (not stock photo or unrelated)?
3. Does it match the description provided?
4. Estimate realistic hours for this task (0.5 to 8 hours)

Be supportive but honest. Confidence score should be:
- 0.9-1.0: Clear evidence of work
- 0.7-0.89: Likely valid with minor uncertainties
- 0.5-0.69: Uncertain, needs review
- Below 0.5: Likely invalid or unclear`, category, description)
}

func (g *GeminiService) parseVerificationResponse(responseText string) (*VerificationResult, error) {
	// Clean the response - remove markdown code blocks if present
	responseText = strings.TrimSpace(responseText)
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var result VerificationResult
	err := json.Unmarshal([]byte(responseText), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w, response: %s", err, responseText)
	}

	// Validate the result
	if result.ConfidenceScore < 0 || result.ConfidenceScore > 1 {
		result.ConfidenceScore = 0.5
	}
	if result.EstimatedHours <= 0 {
		result.EstimatedHours = 1.0
	}

	return &result, nil
}

func (g *GeminiService) Close() {
	if g.client != nil {
		g.client.Close()
	}
}

// CalendarInsights represents the AI-generated insights
type CalendarInsights struct {
	FreeTimeSlots       []FreeTimeSlot       `json:"free_time_slots"`
	WorkPattern         WorkPattern          `json:"work_pattern"`
	BusinessSuggestions []BusinessSuggestion `json:"business_suggestions"`
}

type FreeTimeSlot struct {
	Day       string `json:"day"`
	TimeRange string `json:"time_range"`
	Duration  string `json:"duration"`
}

type WorkPattern struct {
	BusiestDay     string  `json:"busiest_day"`
	AvgHoursPerDay float64 `json:"avg_hours_per_day"`
	PeakTime       string  `json:"peak_time"`
	FreeDays       int     `json:"free_days"`
}

type BusinessSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeSlot    string `json:"time_slot"`
	Duration    string `json:"duration"`
}

// AnalyzeWorkPatterns uses AI to analyze work patterns and generate insights
func (g *GeminiService) AnalyzeWorkPatterns(ctx context.Context, workLogs []models.WorkLog, businessProfile models.BusinessProfile, hasBusinessProfile bool) (*CalendarInsights, error) {
	// Build a summary of work logs
	workSummary := g.buildWorkSummary(workLogs)

	// Create prompt for AI analysis
	prompt := g.createPatternAnalysisPrompt(workSummary, businessProfile, hasBusinessProfile)

	// Get the generative model
	model := g.client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(4096) // Increased to allow full JSON response with suggestions

	// Generate insights
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini AI")
	}

	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Parse the response
	return g.parseCalendarInsights(responseText)
}

func (g *GeminiService) buildWorkSummary(workLogs []models.WorkLog) string {
	if len(workLogs) == 0 {
		return "No work logs available"
	}

	summary := "Work logs from the last 30 days:\n"
	dayMap := make(map[string][]models.WorkLog)

	for _, log := range workLogs {
		day := log.CreatedAt.Weekday().String()
		dayMap[day] = append(dayMap[day], log)
	}

	for day, logs := range dayMap {
		totalHours := 0.0
		categories := make(map[string]int)
		for _, log := range logs {
			totalHours += log.EstimatedHours
			categories[string(log.Category)]++
		}
		summary += fmt.Sprintf("- %s: %.1f hours, %d tasks\n", day, totalHours, len(logs))
	}

	return summary
}

func (g *GeminiService) createPatternAnalysisPrompt(workSummary string, businessProfile models.BusinessProfile, hasBusinessProfile bool) string {
	basePrompt := fmt.Sprintf(`Analyze work patterns and identify free time. Respond with ONLY valid JSON, no extra text.

%s

Required JSON format:
{
  "free_time_slots": [{"day": "Monday", "time_range": "2-4 PM", "duration": "2 hours"}],
  "work_pattern": {"busiest_day": "Monday", "avg_hours_per_day": 3.5, "peak_time": "Morning", "free_days": 2},
  "business_suggestions": []
}

Identify 3-4 free time slots.`, workSummary)

	if hasBusinessProfile {
		basePrompt += fmt.Sprintf(`

Business: %s (%s) targeting %s.
Add 3-4 brief business suggestions (max 80 chars each) in format:
{"title": "...", "description": "...", "time_slot": "Mon 2-4PM", "duration": "2h"}`,
			businessProfile.BusinessName, businessProfile.Industry, businessProfile.TargetAudience)
	} else {
		basePrompt += `

No business profile. Keep business_suggestions as empty array.`
	}

	return basePrompt
}

func (g *GeminiService) parseCalendarInsights(responseText string) (*CalendarInsights, error) {
	// Clean the response
	responseText = strings.TrimSpace(responseText)
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// Try to fix incomplete JSON by adding closing braces if needed
	if !strings.HasSuffix(responseText, "}") {
		// Count opening and closing braces
		openBraces := strings.Count(responseText, "{")
		closeBraces := strings.Count(responseText, "}")
		openBrackets := strings.Count(responseText, "[")
		closeBrackets := strings.Count(responseText, "]")

		// Add missing closing characters
		for i := 0; i < openBrackets-closeBrackets; i++ {
			responseText += "]"
		}
		for i := 0; i < openBraces-closeBraces; i++ {
			responseText += "}"
		}
	}

	var insights CalendarInsights
	err := json.Unmarshal([]byte(responseText), &insights)
	if err != nil {
		// Return a more helpful error with truncated response
		truncated := responseText
		if len(truncated) > 500 {
			truncated = truncated[:500] + "..."
		}
		return nil, fmt.Errorf("failed to parse calendar insights: %w, response (truncated): %s", err, truncated)
	}

	return &insights, nil
}