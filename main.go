package main

import (
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
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
