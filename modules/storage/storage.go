package storage

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"restorapp/modules/auth"
)

var presignClient *s3.PresignClient
var bucketName string
var endpointHost string

func InitStorage() {
	bucketName = os.Getenv("TIGRIS_BUCKET")
	accessKey := os.Getenv("TIGRIS_ACCESS_KEY_ID")
	secretKey := os.Getenv("TIGRIS_ACESSS_KEY_SECRET")
	endpointURL := os.Getenv("TIGRIS_ENDPOINT_URL")
	// Strip scheme for public URL construction (e.g. "https://t3.storage.dev" → "t3.storage.dev")
	endpointHost = strings.TrimPrefix(strings.TrimPrefix(endpointURL, "https://"), "http://")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatal("Failed to load AWS config", "error", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
		o.Region = "auto"
		o.UsePathStyle = false
	})

	presignClient = s3.NewPresignClient(s3Client)
	log.Info("Tigris storage initialized")
}

func GeneratePresignedURL(ctx context.Context, filename string, folder string) (string, string, error) {
	ext := filepath.Ext(filename)
	if folder == "" {
		folder = "products"
	}
	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return "", "", err
	}

	publicURL := fmt.Sprintf("https://%s.%s/%s", bucketName, endpointHost, key)

	return presignedReq.URL, publicURL, nil
}

func StorageController(router *gin.Engine) {
	upload := router.Group("/upload")
	upload.Use(auth.AuthMiddleware())
	upload.Use(auth.EmailVerifiedMiddleware())
	upload.POST("/presign", presignHandler)
}

func presignHandler(ctx *gin.Context) {
	var req struct {
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		Folder      string `json:"folder"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if !strings.HasPrefix(req.ContentType, "image/") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only image uploads are allowed"})
		return
	}

	presignedURL, publicURL, err := GeneratePresignedURL(ctx, req.Filename, req.Folder)
	if err != nil {
		log.Error("Failed to generate presigned URL", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"presignedUrl": presignedURL,
		"publicUrl":    publicURL,
	})
}
