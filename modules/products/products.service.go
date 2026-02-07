package products

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

func getProductsHandler(ctx *gin.Context) {
	productList := []ProductsWithImagesAndCategories{}

	products, err := db.Queries.GetProducts(ctx)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		return
	}

	images, err := db.Queries.GetProductsImages(ctx)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		return
	}

	categories, err := db.Queries.GetProductsCategories(ctx)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		return
	}

	productsImages := make(map[uuid.UUID][]client.ProductImage)
	productsCategories := make(map[uuid.UUID][]client.GetProductsCategoriesRow)

	for _, image := range images {
		productsImages[image.ProductID] = append(productsImages[image.ProductID], image)
	}

	for _, category := range categories {
		productsCategories[category.ProductID] = append(productsCategories[category.ProductID], category)
	}

	for _, product := range products {
		productToAdd := ProductsWithImagesAndCategories{
			Product:    product,
			Images:     productsImages[product.ID],
			Categories: productsCategories[product.ID],
		}

		productList = append(productList, productToAdd)
	}

	ctx.JSON(http.StatusOK, productList)
}

func createProductHandler(ctx *gin.Context) {
	productToCreate := client.CreateProductParams{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&productToCreate)
	if err != nil {
		log.Error("Unproceasable entity", err)
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	createdProduct, errDB := db.Queries.CreateProduct(ctx, productToCreate)
	if errDB != nil {
		log.Error("Error happen on db", errDB)
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"product": createdProduct,
	})
}

func deleteProductHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")

	productUUID, errUUID := uuid.Parse(productIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.Status(http.StatusBadRequest)
		return
	}

	product, errDB := db.Queries.DeleteProduct(ctx, productUUID)
	if errDB != nil {
		log.Error("Error on db", errDB)
		ctx.Status(http.StatusBadRequest)
		return
	}

	if len(product) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}

func updateProductHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, errUUID := uuid.Parse(productIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.Status(http.StatusBadRequest)
		return
	}
	productToUpdate := client.UpdateProductParams{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&productToUpdate)
	if err != nil {
		log.Error("Bad Request", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	productToUpdate.ID = productUUID

	product, errDB := db.Queries.UpdateProduct(ctx, productToUpdate)
	if errDB != nil {
		log.Error("Error happen on db", errDB)
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	if len(product) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "updated successfully",
	})
}

func getProductByIdHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, errUUID := uuid.Parse(productIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.Status(http.StatusBadRequest)
		return
	}
	product, err := db.Queries.GetProductById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	category, err := db.Queries.GetProductCategoriesById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	images, err := db.Queries.GetProductImagesById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	productData := ProductWithImagesAndCategories{
		Product:    product,
		Categories: category,
		Images:     images,
	}

	ctx.JSON(http.StatusOK, productData)
}
