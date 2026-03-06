package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sheleads-backend/config"
	"time"
)

const (
	runwayAPIBase = "https://api.dev.runwayml.com/v1"
)

// RunwayService handles Runway AI image generation
type RunwayService struct {
	apiKey     string
	httpClient *http.Client
}

// RunwayImageRequest represents the request to generate an image
type RunwayImageRequest struct {
	PromptText      string                   `json:"promptText"`
	Ratio           string                   `json:"ratio"`
	Seed            int64                    `json:"seed"`
	Model           string                   `json:"model"`
	ReferenceImages []RunwayReferenceImage   `json:"referenceImages,omitempty"`
}

// RunwayReferenceImage represents a reference image for generation
type RunwayReferenceImage struct {
	URI string `json:"uri"`
}

// RunwayTaskResponse represents the initial task creation response
type RunwayTaskResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// RunwayTaskStatus represents the task status response
type RunwayTaskStatus struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Output []string `json:"output"`
}

// NewRunwayService creates a new Runway service instance
func NewRunwayService() (*RunwayService, error) {
	apiKey := config.AppConfig.RunwayAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("Runway API key not configured")
	}

	return &RunwayService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}, nil
}

// GenerateImage generates an image using Runway AI
func (r *RunwayService) GenerateImage(ctx context.Context, prompt string, referenceImageURLs []string) (string, error) {
	// Prepare request
	request := RunwayImageRequest{
		PromptText: prompt,
		Ratio:      "1920:1080",
		Seed:       time.Now().Unix() % 4294967295,
		Model:      "gen4_image",
	}

	// Add reference images if provided
	if len(referenceImageURLs) > 0 {
		for _, url := range referenceImageURLs {
			if url != "" {
				request.ReferenceImages = append(request.ReferenceImages, RunwayReferenceImage{
					URI: url,
				})
			}
		}
	}

	// Create task
	taskID, err := r.createImageTask(request)
	if err != nil {
		return "", fmt.Errorf("failed to create Runway task: %w", err)
	}

	// Poll for completion
	imageURL, err := r.waitForTaskCompletion(taskID, 60) // 60 attempts, 5 seconds each = 5 minutes max
	if err != nil {
		return "", fmt.Errorf("failed to complete Runway task: %w", err)
	}

	// Download image to local storage
	localPath, err := r.downloadImage(imageURL, taskID)
	if err != nil {
		return "", fmt.Errorf("failed to download generated image: %w", err)
	}

	return localPath, nil
}

// createImageTask creates a new image generation task
func (r *RunwayService) createImageTask(request RunwayImageRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", runwayAPIBase+"/text_to_image", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Runway-Version", "2024-11-06")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("runway API error: %s - %s", resp.Status, string(body))
	}

	var taskResp RunwayTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return "", err
	}

	return taskResp.ID, nil
}

// waitForTaskCompletion polls the task status until completion
func (r *RunwayService) waitForTaskCompletion(taskID string, maxAttempts int) (string, error) {
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(5 * time.Second)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks/%s", runwayAPIBase, taskID), nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+r.apiKey)
		req.Header.Set("X-Runway-Version", "2024-11-06")

		resp, err := r.httpClient.Do(req)
		if err != nil {
			continue
		}

		var status RunwayTaskStatus
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if status.Status == "SUCCEEDED" && len(status.Output) > 0 {
			return status.Output[0], nil
		} else if status.Status == "FAILED" || status.Status == "CANCELLED" {
			return "", fmt.Errorf("runway task failed with status: %s", status.Status)
		}
	}

	return "", fmt.Errorf("runway task timed out after %d attempts", maxAttempts)
}

// downloadImage downloads the generated image to local storage
func (r *RunwayService) downloadImage(imageURL, taskID string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create directory if not exists
	dir := "generated_images"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Save file
	localPath := filepath.Join(dir, fmt.Sprintf("%s.jpg", taskID))
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return localPath, nil
}

