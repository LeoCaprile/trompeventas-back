package products

import (
	"restorapp/modules/auth"

	"github.com/gin-gonic/gin"
)

func ProductsController(router *gin.Engine) {
	router.GET("/products", getProductsHandler)
	router.GET("/products/:id", getProductByIdHandler)

	products := router.Group("/products")
	products.Use(auth.AuthMiddleware())
	products.POST("/", createProductHandler)
	products.GET("/me", getMyProductsHandler)
	products.DELETE("/me/:id", deleteMyProductHandler)
	products.PUT("/me/:id", updateMyProductHandler)

	publish := router.Group("/products")
	publish.Use(auth.AuthMiddleware())
	publish.Use(auth.EmailVerifiedMiddleware())
	publish.POST("/publish", publishProductHandler)
}
