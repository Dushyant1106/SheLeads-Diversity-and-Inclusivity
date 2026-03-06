package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sheleads-backend/config"
	"sheleads-backend/models"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ContentAIService handles AI-powered content generation
type ContentAIService struct {
	client *genai.Client
}

// NewContentAIService creates a new content AI service
func NewContentAIService() *ContentAIService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.AppConfig.GeminiAPIKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		return nil
	}

	return &ContentAIService{
		client: client,
	}
}

// GenerateBlogPost generates a blog post using AI
func (s *ContentAIService) GenerateBlogPost(ctx context.Context, business *models.BusinessProfile, request *models.BlogGenerationRequest) (string, string, error) {
	model := s.client.GenerativeModel("gemini-2.5-flash")
	
	// Build the prompt
	tone := request.Tone
	if tone == "" {
		tone = "professional and engaging"
	}
	
	length := request.Length
	if length == "" {
		length = "medium (800-1200 words)"
	}
	
	keywords := strings.Join(request.Keywords, ", ")
	
	prompt := fmt.Sprintf(`You are a professional content writer helping women entrepreneurs create engaging blog content.

Business Information:
- Business Name: %s
- Industry: %s
- Location: %s
- Description: %s
- Target Audience: %s
- Unique Selling Points: %s

Task: Write a complete blog post about "%s"

Requirements:
- Tone: %s
- Length: %s
- Include these keywords naturally: %s
- Make it SEO-friendly with proper headings (H2, H3)
- Include an engaging introduction and strong conclusion
- Add actionable tips or insights
- Make it relevant to the business and target audience

Return ONLY a JSON object with this exact structure:
{
  "title": "Blog post title here",
  "content": "Full HTML blog content with <h2>, <h3>, <p>, <ul>, <li> tags"
}`,
		business.BusinessName,
		business.Industry,
		business.Location,
		business.Description,
		business.TargetAudience,
		strings.Join(business.UniqueSellingPoints, ", "),
		request.Topic,
		tone,
		length,
		keywords,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate blog content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", "", fmt.Errorf("no content generated")
	}

	// Extract the text response
	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	
	// Clean up the response (remove markdown code blocks if present)
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// Parse JSON response
	var result struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		// If JSON parsing fails, return the raw content
		log.Printf("Failed to parse JSON response: %v", err)
		return "Generated Blog Post", responseText, nil
	}

	return result.Title, result.Content, nil
}

// GenerateSocialPost generates social media posts using AI
func (s *ContentAIService) GenerateSocialPost(ctx context.Context, business *models.BusinessProfile, request *models.SocialPostGenerationRequest) (map[string]string, error) {
	model := s.client.GenerativeModel("gemini-2.5-flash")
	
	tone := request.Tone
	if tone == "" {
		tone = "engaging and authentic"
	}
	
	platforms := strings.Join(request.Platforms, ", ")
	hashtags := ""
	if len(request.Hashtags) > 0 {
		hashtags = strings.Join(request.Hashtags, ", ")
	}

	prompt := fmt.Sprintf(`You are a social media expert helping women entrepreneurs create engaging posts.

Business Information:
- Business Name: %s
- Industry: %s
- Target Audience: %s
- Unique Selling Points: %s

Task: Create social media posts about "%s" for these platforms: %s

Requirements:
- Tone: %s
- Tailor each post to the platform's best practices
- LinkedIn: Professional, 150-200 words
- Twitter: Concise, max 280 characters
- Instagram: Visual-focused, engaging caption with emojis
- Facebook: Conversational, 100-150 words
`, business.BusinessName, business.Industry, business.TargetAudience, 
   strings.Join(business.UniqueSellingPoints, ", "), request.Topic, platforms, tone)

	if hashtags != "" {
		prompt += fmt.Sprintf("- Include these hashtags: %s\n", hashtags)
	}

	prompt += `
Return ONLY a JSON object with platform names as keys and post content as values.
Example: {"linkedin": "post content", "twitter": "post content"}`

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate social posts: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var posts map[string]string
	if err := json.Unmarshal([]byte(responseText), &posts); err != nil {
		return nil, fmt.Errorf("failed to parse social posts: %w", err)
	}

	return posts, nil
}

