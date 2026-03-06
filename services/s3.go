package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sheleads-backend/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Service handles AWS S3 operations
type S3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new S3 service instance
func NewS3Service() (*S3Service, error) {
	accessKey := config.AppConfig.AWSAccessKeyID
	secretKey := config.AppConfig.AWSSecretAccessKey
	bucketName := config.AppConfig.AWSS3BucketName
	region := config.AppConfig.AWSRegion

	if accessKey == "" || secretKey == "" || bucketName == "" {
		return nil, fmt.Errorf("AWS credentials not configured")
	}

	if region == "" {
		region = "us-east-1" // Default region
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Service{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads a file to S3 and returns the public URL
func (s *S3Service) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, userID, assetType string) (string, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique S3 key
	uniqueID := uuid.New().String()[:8]
	timestamp := time.Now().Format("20060102")
	fileExt := filepath.Ext(fileHeader.Filename)
	s3Key := fmt.Sprintf("user-assets/%s/%s/%s_%s%s", userID, assetType, timestamp, uniqueID, fileExt)

	// Determine content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentType(fileExt)
	}

	// Upload to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate public URL
	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, s3Key)
	return s3URL, nil
}

// UploadFileFromPath uploads a file from local path to S3
func (s *S3Service) UploadFileFromPath(ctx context.Context, filePath, s3Key, contentType string) (string, error) {
	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate public URL
	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, s3Key)
	return s3URL, nil
}

// UploadBlog uploads a blog HTML file to S3
func (s *S3Service) UploadBlog(ctx context.Context, htmlContent, contentID string) (string, error) {
	s3Key := fmt.Sprintf("blogs/%s.html", contentID)
	
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader([]byte(htmlContent)),
		ContentType: aws.String("text/html"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload blog to S3: %w", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, s3Key)
	return s3URL, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(ctx context.Context, s3Key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

// getContentType returns the content type based on file extension
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

