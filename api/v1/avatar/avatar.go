package avatar

import (
	"net/http"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/avatar")
	{
		g.POST("/", createAvatar)
		g.GET("/:id", getAvatar)
		g.GET("/", getAllAvatars)
		g.PUT("/:id", updateAvatar)
		g.DELETE("/:id", deleteAvatar)
	}
}

// Create Avatar
func createAvatar(c *gin.Context) {
	var avatar models.Avatar
	if err := c.ShouldBindJSON(&avatar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := dbconfig.GetDb()
	if err := db.Create(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, avatar)
}

// Get Avatar by ID
func getAvatar(c *gin.Context) {
	id := c.Param("id")
	var avatar models.Avatar

	db := dbconfig.GetDb()
	if err := db.First(&avatar, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	c.JSON(http.StatusOK, avatar)
}

// Get all Avatars
func getAllAvatars(c *gin.Context) {
	var avatars []models.Avatar

	db := dbconfig.GetDb()
	if err := db.Find(&avatars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avatars)
}

// Update Avatar
func updateAvatar(c *gin.Context) {
	id := c.Param("id")
	var avatar models.Avatar

	db := dbconfig.GetDb()
	if err := db.First(&avatar, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	if err := c.ShouldBindJSON(&avatar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avatar)
}

// Delete Avatar
func deleteAvatar(c *gin.Context) {
	id := c.Param("id")
	var avatar models.Avatar

	db := dbconfig.GetDb()
	if err := db.First(&avatar, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	if err := db.Delete(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar deleted successfully"})
}

