package locations

import (
	"embed"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed comunas-regiones.json
var dataFS embed.FS

type regionEntry struct {
	Region  string   `json:"region"`
	Comunas []string `json:"comunas"`
}

type regionsData struct {
	Regiones []regionEntry `json:"regiones"`
}

var cachedData *regionsData

func loadData() (*regionsData, error) {
	if cachedData != nil {
		return cachedData, nil
	}

	raw, err := dataFS.ReadFile("comunas-regiones.json")
	if err != nil {
		return nil, err
	}

	var d regionsData
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}

	cachedData = &d
	return cachedData, nil
}

func LocationsController(router *gin.Engine) {
	router.GET("/locations/regions", handleGetRegions)
}

func handleGetRegions(c *gin.Context) {
	data, err := loadData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load regions data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
