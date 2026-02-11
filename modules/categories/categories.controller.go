package categories

import (
	"restorapp/modules/auth"

	"github.com/gin-gonic/gin"
)

func CategoriesController(router *gin.Engine) {
	router.GET("/categories", getCategoriesHandler)

	categoriesRouter := router.Group("/categories")
	categoriesRouter.Use(auth.AuthMiddleware())
	categoriesRouter.POST("/", createCategoriesHandler)
	categoriesRouter.DELETE("/:id", deleteCategoriesHandler)
	categoriesRouter.PUT("/:id", updateCategoriesHandler)
}
