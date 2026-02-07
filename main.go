package main

import (
	"restorapp/db"
	"restorapp/modules/auth"
	"restorapp/modules/categories"
	"restorapp/modules/email"
	"restorapp/modules/products"

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

	auth.InitAuth(router)

	products.ProductsController(router)
	categories.CategoriesController(router)

	router.Run()
}
