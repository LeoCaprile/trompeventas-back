package categories

import (
	"restorapp/modules/auth"

	"github.com/gin-gonic/gin"
)

func CategoriesController(router *gin.Engine) {
	categoriesRouter := router.Group("/categories")

	categoriesRouter.Use(auth.AuthMiddleware())
	categoriesRouter.GET("/", getCategoriesHandler)
	categoriesRouter.POST("/", createCategoriesHandler)
	categoriesRouter.DELETE("/:id", deleteCategoriesHandler)
	categoriesRouter.POST("/:id", updateCategoriesHandler)
}
