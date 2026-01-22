package categories

import "github.com/gin-gonic/gin"

func CategoriesController(router *gin.Engine) {
	router.GET("/categories", getCategoriesHandler)
	router.POST("/categories", createCategoriesHandler)
	router.DELETE("/categories/:id", deleteCategoriesHandler)
	router.POST("/categories/:id", updateCategoriesHandler)
}
