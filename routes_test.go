package main

import (
	"testing"

	"restorapp/modules/categories"
	"restorapp/modules/comments"
	"restorapp/modules/locations"
	"restorapp/modules/products"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	products.ProductsController(router)
	categories.CategoriesController(router)
	comments.CommentsController(router)
	locations.LocationsController(router)
	return router
}

type routeEntry struct {
	method string
	path   string
}

func TestProductRoutes(t *testing.T) {
	router := setupRouter()
	routes := router.Routes()
	routeSet := make(map[routeEntry]bool)
	for _, r := range routes {
		routeSet[routeEntry{r.Method, r.Path}] = true
	}

	expected := []routeEntry{
		{"GET", "/products"},
		{"GET", "/products/:id"},
		{"POST", "/products/"},
		{"GET", "/products/me"},
		{"DELETE", "/products/me/:id"},
		{"PUT", "/products/me/:id"},
		{"POST", "/products/publish"},
	}

	for _, e := range expected {
		if !routeSet[e] {
			t.Errorf("expected route %s %s not found", e.method, e.path)
		}
	}

	// Verify old incorrect routes are NOT registered
	forbidden := []routeEntry{
		{"POST", "/products/me/:id"},   // should be PUT, not POST
		{"DELETE", "/products/:id"},     // unprotected admin route removed
		{"POST", "/products/:id"},      // unprotected admin route removed
	}

	for _, f := range forbidden {
		if routeSet[f] {
			t.Errorf("route %s %s should NOT be registered", f.method, f.path)
		}
	}
}

func TestCategoryRoutes(t *testing.T) {
	router := setupRouter()
	routes := router.Routes()
	routeSet := make(map[routeEntry]bool)
	for _, r := range routes {
		routeSet[routeEntry{r.Method, r.Path}] = true
	}

	expected := []routeEntry{
		{"GET", "/categories"},
		{"POST", "/categories/"},
		{"DELETE", "/categories/:id"},
		{"PUT", "/categories/:id"},
	}

	for _, e := range expected {
		if !routeSet[e] {
			t.Errorf("expected route %s %s not found", e.method, e.path)
		}
	}

	// Verify old POST update route is not registered
	if routeSet[routeEntry{"POST", "/categories/:id"}] {
		t.Error("route POST /categories/:id should NOT be registered (should be PUT)")
	}
}

func TestCommentRoutes(t *testing.T) {
	router := setupRouter()
	routes := router.Routes()
	routeSet := make(map[routeEntry]bool)
	for _, r := range routes {
		routeSet[routeEntry{r.Method, r.Path}] = true
	}

	expected := []routeEntry{
		{"GET", "/products/:id/comments"},
		{"POST", "/products/:id/comments/"},
		{"DELETE", "/comments/:id"},
		{"PUT", "/comments/:id/vote"},
		{"DELETE", "/comments/:id/vote"},
	}

	for _, e := range expected {
		if !routeSet[e] {
			t.Errorf("expected route %s %s not found", e.method, e.path)
		}
	}

	// Verify vote uses PUT not POST
	if routeSet[routeEntry{"POST", "/comments/:id/vote"}] {
		t.Error("route POST /comments/:id/vote should NOT be registered (should be PUT)")
	}
}

func TestLocationRoutes(t *testing.T) {
	router := setupRouter()
	routes := router.Routes()
	routeSet := make(map[routeEntry]bool)
	for _, r := range routes {
		routeSet[routeEntry{r.Method, r.Path}] = true
	}

	expected := []routeEntry{
		{"GET", "/locations/regions"},
	}

	for _, e := range expected {
		if !routeSet[e] {
			t.Errorf("expected route %s %s not found", e.method, e.path)
		}
	}
}
