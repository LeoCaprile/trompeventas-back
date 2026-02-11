package products

import (
	"restorapp/db/client"
)

type ProductsWithImagesAndCategories struct {
	Product    client.Product                    `json:"product"`
	Images     []client.ProductImage             `json:"images"`
	Categories []client.GetProductsCategoriesRow `json:"categories"`
}

type SellerInfo struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Region string `json:"region"`
	City   string `json:"city"`
}

type ProductWithImagesAndCategories struct {
	Product    client.Product                       `json:"product"`
	Images     []client.ProductImage                `json:"images"`
	Categories []client.GetProductCategoriesByIdRow `json:"categories"`
	Seller     *SellerInfo                          `json:"seller"`
}
