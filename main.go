package main

import (
	"os"
	"strings"
	"restorapp/db"
	"restorapp/modules/auth"
	"restorapp/modules/categories"
	"restorapp/modules/email"
	"restorapp/modules/locations"
	"restorapp/modules/comments"
	"restorapp/modules/products"
	"restorapp/modules/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Get allowed origins from environment variable (required)
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		panic("ALLOWED_ORIGINS environment variable is required")
	}

	origins := strings.Split(allowedOrigins, ",")
	// Trim whitespace from each origin
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Cookie", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	conn := db.InitDBClient()
	defer conn.Close()
	email.InitResendClient()
	storage.InitStorage()

	auth.InitAuth(router)

	products.ProductsController(router)
	categories.CategoriesController(router)
	comments.CommentsController(router)
	locations.LocationsController(router)
	storage.StorageController(router)

	router.Run()
}
