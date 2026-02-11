package categories

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"restorapp/db"
	"restorapp/db/client"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func getCategoriesHandler(ctx *gin.Context) {
	data, err := db.Queries.GetCategories(ctx)
	if err != nil {
		log.Error("Could not retrieve data of categories", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func createCategoriesHandler(ctx *gin.Context) {
	category := struct {
		Name string `json:"name"`
	}{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&category)
	if err != nil {
		log.Error("Unprocessable entity", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body"})
		return
	}

	createdCategories, errDB := db.Queries.CreateCategory(ctx, category.Name)
	if errDB != nil {
		log.Error("Error creating category in db", errDB)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to create category"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": createdCategories,
	})
}

func deleteCategoriesHandler(ctx *gin.Context) {
	categoryIdParam := ctx.Param("id")

	categoryUUID, errUUID := uuid.Parse(categoryIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, errDB := db.Queries.DeleteCategory(ctx, categoryUUID)
	if errDB != nil {
		log.Error("Error deleting category", errDB)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	if len(category) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}

func updateCategoriesHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, errUUID := uuid.Parse(productIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	categoryToUpdate := client.UpdateCategoryParams{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&categoryToUpdate)
	if err != nil {
		log.Error("Bad Request", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	categoryToUpdate.ID = productUUID

	categories, errDB := db.Queries.UpdateCategory(ctx, categoryToUpdate)
	if errDB != nil {
		log.Error("Error updating category", errDB)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	if len(categories) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
	})
}
