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
	data, err := db.Queries.GetProducts(ctx)
	if err != nil {
		log.Error("Could not retrieve data of tables", err)
		return
	}

	ctx.JSON(http.StatusOK, data)
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

	errDB := db.Queries.DeleteProduct(ctx, productUUID)
	if errDB != nil {
		log.Error("Error on db", errDB)
		ctx.Status(http.StatusBadRequest)
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
	product := client.UpdateProductParams{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&product)
	if err != nil {
		log.Error("Bad Request", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	product.ID = productUUID

	errDB := db.Queries.UpdateProduct(ctx, product)
	if errDB != nil {
		log.Error("Error happen on db", errDB)
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "updated successfully",
	})
}
