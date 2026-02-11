package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"restorapp/db"
	"restorapp/db/client"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func getProductsHandler(ctx *gin.Context) {
	productList := []ProductsWithImagesAndCategories{}

	q := ctx.Query("q")
	var products []client.Product
	var err error

	if q != "" {
		products, err = db.Queries.SearchProducts(ctx, pgtype.Text{String: q, Valid: true})
	} else {
		products, err = db.Queries.GetProducts(ctx)
	}
	if err != nil {
		log.Error("Could not retrieve data of products", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	images, err := db.Queries.GetProductsImages(ctx)
	if err != nil {
		log.Error("Could not retrieve data of product images", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product images"})
		return
	}

	categories, err := db.Queries.GetProductsCategories(ctx)
	if err != nil {
		log.Error("Could not retrieve data of product categories", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product categories"})
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
		log.Error("Unprocessable entity", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body"})
		return
	}

	createdProduct, errDB := db.Queries.CreateProduct(ctx, productToCreate)
	if errDB != nil {
		log.Error("Error creating product in db", errDB)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"product": createdProduct,
	})
}

func getProductByIdHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, errUUID := uuid.Parse(productIdParam)
	if errUUID != nil {
		log.Error("Error parsing UUID", errUUID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	product, err := db.Queries.GetProductById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve product", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	category, err := db.Queries.GetProductCategoriesById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve product categories", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product categories"})
		return
	}

	images, err := db.Queries.GetProductImagesById(ctx, productUUID)
	if err != nil {
		log.Error("Could not retrieve product images", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product images"})
		return
	}

	var seller *SellerInfo
	if product.UserID.Valid {
		user, err := db.Queries.GetUserById(ctx, product.UserID.Bytes)
		if err == nil {
			seller = &SellerInfo{
				Name:   user.Name,
				Image:  user.Image.String,
				Region: user.Region.String,
				City:   user.City.String,
			}
		}
	}

	productData := ProductWithImagesAndCategories{
		Product:    product,
		Categories: category,
		Images:     images,
		Seller:     seller,
	}

	ctx.JSON(http.StatusOK, productData)
}

func getMyProductsHandler(ctx *gin.Context) {
	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	products, err := db.Queries.GetProductsByUserId(ctx, pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil {
		log.Error("Could not retrieve user products", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	productList := []ProductsWithImagesAndCategories{}
	for _, product := range products {
		images, _ := db.Queries.GetProductImagesById(ctx, product.ID)
		categories, _ := db.Queries.GetProductCategoriesById(ctx, product.ID)

		catRows := []client.GetProductsCategoriesRow{}
		for _, c := range categories {
			catRows = append(catRows, client.GetProductsCategoriesRow{
				ID:        c.ID,
				ProductID: c.ProductID,
				Name:      c.Name,
			})
		}

		productList = append(productList, ProductsWithImagesAndCategories{
			Product:    product,
			Images:     images,
			Categories: catRows,
		})
	}

	ctx.JSON(http.StatusOK, productList)
}

func deleteMyProductHandler(ctx *gin.Context) {
	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	productIdParam := ctx.Param("id")
	productUUID, err := uuid.Parse(productIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = db.Queries.DeleteProductByOwner(ctx, client.DeleteProductByOwnerParams{
		ID:     productUUID,
		UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
	})
	if err != nil {
		log.Error("Error deleting product", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func updateMyProductHandler(ctx *gin.Context) {
	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	productIdParam := ctx.Param("id")
	productUUID, err := uuid.Parse(productIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Verify ownership
	product, err := db.Queries.GetProductById(ctx, productUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if product.UserID.Bytes != userUUID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Not your product"})
		return
	}

	productToUpdate := client.UpdateProductParams{}
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(&productToUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	productToUpdate.ID = productUUID

	updated, err := db.Queries.UpdateProduct(ctx, productToUpdate)
	if err != nil {
		log.Error("Error updating product", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	if len(updated) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": updated[0]})
}

type PublishProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Condition   string   `json:"condition"`
	Negotiable  string   `json:"negotiable"`
	Categories  []string `json:"categories"`
	ImageUrls   []string `json:"imageUrls"`
}

func publishProductHandler(ctx *gin.Context) {
	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req PublishProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if req.Price <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}
	if len(req.ImageUrls) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		log.Error("Failed to begin transaction", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer tx.Rollback(context.Background())

	qtx := db.Queries.WithTx(tx)

	condition := req.Condition
	if condition == "" {
		condition = "Nuevo"
	}
	negotiable := req.Negotiable
	if negotiable == "" {
		negotiable = "No conversable"
	}

	product, err := qtx.CreateProduct(ctx, client.CreateProductParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Price:       req.Price,
		UserID:      pgtype.UUID{Bytes: userUUID, Valid: true},
		Condition:   condition,
		State:       "Disponible",
		Negotiable:  negotiable,
	})
	if err != nil {
		log.Error("Failed to create product", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	var images []client.ProductImage
	for _, imageUrl := range req.ImageUrls {
		image, err := qtx.CreateProductImage(ctx, client.CreateProductImageParams{
			ProductID: product.ID,
			ImageUrl:  imageUrl,
		})
		if err != nil {
			log.Error("Failed to create product image", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product image"})
			return
		}
		images = append(images, image)
	}

	var categories []client.ProductsCategory
	for _, catIdStr := range req.Categories {
		catUUID, err := uuid.Parse(catIdStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid category ID: %s", catIdStr)})
			return
		}
		cat, err := qtx.CreateProductCategory(ctx, client.CreateProductCategoryParams{
			ProductID:  product.ID,
			CategoryID: catUUID,
		})
		if err != nil {
			log.Error("Failed to create product category", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product category"})
			return
		}
		categories = append(categories, cat)
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Error("Failed to commit transaction", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product published successfully",
		"product": product,
		"images":  images,
	})
}
