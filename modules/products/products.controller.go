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
	products.DELETE("/:id", deleteProductHandler)
	products.POST("/:id", updateProductHandler)
}
