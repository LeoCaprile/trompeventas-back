package products

import (
	"restorapp/db/client"
)

type ProductsWithImagesAndCategories struct {
	Product    client.Product                    `json:"product"`
	Images     []client.ProductImage             `json:"images"`
	Categories []client.GetProductsCategoriesRow `json:"categories"`
}

type ProductWithImagesAndCategories struct {
	Product    client.Product                       `json:"product"`
	Images     []client.ProductImage                `json:"images"`
	Categories []client.GetProductCategoriesByIdRow `json:"categories"`
}
