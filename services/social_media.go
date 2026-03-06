package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sheleads-backend/config"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// SocialMediaService handles posting to social media platforms
type SocialMediaService struct {
	twitterClient *twitter.Client
}

// NewSocialMediaService creates a new social media service instance
func NewSocialMediaService() (*SocialMediaService, error) {
	service := &SocialMediaService{}

	// Initialize Twitter client if credentials are available
	if config.AppConfig.TwitterAPIKey != "" &&
		config.AppConfig.TwitterAPISecret != "" &&
		config.AppConfig.TwitterAccessToken != "" &&
		config.AppConfig.TwitterAccessTokenSecret != "" {
		
		oauthConfig := oauth1.NewConfig(
			config.AppConfig.TwitterAPIKey,
			config.AppConfig.TwitterAPISecret,
		)
		token := oauth1.NewToken(
			config.AppConfig.TwitterAccessToken,
			config.AppConfig.TwitterAccessTokenSecret,
		)
		httpClient := oauthConfig.Client(oauth1.NoContext, token)
		service.twitterClient = twitter.NewClient(httpClient)
	}

	return service, nil
}

// PostToTwitter posts text and image to Twitter
func (s *SocialMediaService) PostToTwitter(text, imagePath string) (string, error) {
	if s.twitterClient == nil {
		return "", fmt.Errorf("Twitter client not initialized - check credentials")
	}

	// Note: The go-twitter library doesn't support media upload directly
	// For full Twitter API v1.1 support with media upload, you would need to use
	// a different library or implement the media upload endpoint manually
	// For now, we'll post text-only tweets

	// Create tweet (text only for now)
	tweet, _, err := s.twitterClient.Statuses.Update(text, nil)
	if err != nil {
		return "", fmt.Errorf("failed to post tweet: %w", err)
	}

	// Generate tweet URL
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%d", tweet.User.ScreenName, tweet.ID)
	return tweetURL, nil
}

// PostToInstagram posts to Instagram using Python instagrapi library
func (s *SocialMediaService) PostToInstagram(text, imagePath string) (string, error) {
	username := config.AppConfig.InstagramUsername
	password := config.AppConfig.InstagramPassword

	if username == "" || password == "" {
		return "", fmt.Errorf("Instagram credentials not configured")
	}

	if imagePath == "" {
		return "", fmt.Errorf("Instagram requires an image to post")
	}

	// Check if image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("image file not found: %s", imagePath)
	}

	// Create Python script for Instagram posting
	pythonScript := `
import sys
import json
from instagrapi import Client

def post_to_instagram(username, password, image_path, caption):
    try:
        cl = Client()

        # Try to load session if exists
        session_file = "instagram_session.json"
        try:
            cl.load_settings(session_file)
            cl.login(username, password)
        except:
            cl.login(username, password)
            cl.dump_settings(session_file)

        # Upload photo
        media = cl.photo_upload(path=image_path, caption=caption)

        # Return media code for URL construction
        return {
            "success": True,
            "media_code": media.code,
            "media_id": media.id
        }
    except Exception as e:
        return {
            "success": False,
            "error": str(e)
        }

if __name__ == "__main__":
    username = sys.argv[1]
    password = sys.argv[2]
    image_path = sys.argv[3]
    caption = sys.argv[4]

    result = post_to_instagram(username, password, image_path, caption)
    print(json.dumps(result))
`

	// Write Python script to temporary file
	scriptPath := filepath.Join(os.TempDir(), "instagram_post.py")
	if err := os.WriteFile(scriptPath, []byte(pythonScript), 0644); err != nil {
		return "", fmt.Errorf("failed to create Python script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Execute Python script
	cmd := exec.Command("python3", scriptPath, username, password, imagePath, text)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute Instagram posting script: %w, stderr: %s", err, stderr.String())
	}

	// Parse result
	var result struct {
		Success   bool   `json:"success"`
		MediaCode string `json:"media_code"`
		MediaID   string `json:"media_id"`
		Error     string `json:"error"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return "", fmt.Errorf("failed to parse Instagram response: %w, output: %s", err, stdout.String())
	}

	if !result.Success {
		return "", fmt.Errorf("Instagram posting failed: %s", result.Error)
	}

	// Construct Instagram post URL
	postURL := fmt.Sprintf("https://www.instagram.com/p/%s/", result.MediaCode)
	return postURL, nil
}

// PostToFacebook posts to Facebook (placeholder)
func (s *SocialMediaService) PostToFacebook(text, imagePath string) (string, error) {
	// Facebook posting requires Facebook Graph API
	// Similar to Instagram, this would need proper API integration
	
	return "", fmt.Errorf("Facebook posting not yet implemented - requires Facebook Graph API setup")
}

// PostToLinkedIn posts to LinkedIn (placeholder)
func (s *SocialMediaService) PostToLinkedIn(text, imagePath string) (string, error) {
	// LinkedIn posting requires LinkedIn API
	// This would need proper OAuth2 flow and API integration
	
	return "", fmt.Errorf("LinkedIn posting not yet implemented - requires LinkedIn API setup")
}

// PostToSocial posts to the specified platform
func (s *SocialMediaService) PostToSocial(platform, text, imagePath string) (string, error) {
	switch platform {
	case "twitter":
		return s.PostToTwitter(text, imagePath)
	case "instagram":
		return s.PostToInstagram(text, imagePath)
	case "facebook":
		return s.PostToFacebook(text, imagePath)
	case "linkedin":
		return s.PostToLinkedIn(text, imagePath)
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}

