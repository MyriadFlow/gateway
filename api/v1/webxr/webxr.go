package webxr

import (
	"net/http"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/webxr")
	{
		g.POST("/", createWebXR)
		g.GET("/:id", getWebXR)
		g.GET("/", getAllWebXRs)
		g.PUT("/:id", updateWebXR)
		g.DELETE("/:id", deleteWebXR)
	}
}

// Create WebXR
func createWebXR(c *gin.Context) {
	var webxr models.WebXR
	if err := c.ShouldBindJSON(&webxr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := dbconfig.GetDb()
	if err := db.Create(&webxr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, webxr)
}

// Get WebXR by ID
func getWebXR(c *gin.Context) {
	id := c.Param("id")
	var webxr models.WebXR

	db := dbconfig.GetDb()
	if err := db.First(&webxr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WebXR not found"})
		return
	}

	c.JSON(http.StatusOK, webxr)
}

// Get all WebXRs
func getAllWebXRs(c *gin.Context) {
	var webxrs []models.WebXR

	db := dbconfig.GetDb()
	if err := db.Find(&webxrs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webxrs)
}

// Update WebXR
func updateWebXR(c *gin.Context) {
	id := c.Param("id")
	var webxr models.WebXR

	db := dbconfig.GetDb()
	if err := db.First(&webxr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WebXR not found"})
		return
	}

	if err := c.ShouldBindJSON(&webxr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&webxr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webxr)
}

// Delete WebXR
func deleteWebXR(c *gin.Context) {
	id := c.Param("id")
	var webxr models.WebXR

	db := dbconfig.GetDb()
	if err := db.First(&webxr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WebXR not found"})
		return
	}

	if err := db.Delete(&webxr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WebXR deleted successfully"})
}

