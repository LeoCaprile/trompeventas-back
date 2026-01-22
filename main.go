package main

import (
	"context"

	"restorapp/db"
	"restorapp/modules/products"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	conn := db.InitDBClient()
	defer conn.Close(context.Background())

	products.ProductsController(router)

	router.Run()
}
