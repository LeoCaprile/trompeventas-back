package products

import "github.com/gin-gonic/gin"

func ProductsController(router *gin.Engine) {
	router.GET("/products", getProductsHandler)
	router.POST("/products", createProductHandler)
	router.DELETE("/products/:id", deleteProductHandler)
	router.POST("/products/:id", updateProductHandler)
}
