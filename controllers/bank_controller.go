package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBank(c *gin.Context) {
	var bank models.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.GetDB().Create(&bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bank)
}
func GetBank(c *gin.Context) {
	id := c.Param("id")

	var bank models.Bank

	if err := config.GetDB().
		Preload("Branches").
		First(&bank, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}

	c.JSON(http.StatusOK, bank)
}
func GetAllBanks(c *gin.Context) {
	var banks []models.Bank

	if err := config.GetDB().
		Preload("Branches").
		Find(&banks).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, banks)
}
func UpdateBank(c *gin.Context) {
	id := c.Param("id")

	db := config.GetDB()
	var bank models.Bank
	if err := db.First(&bank, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}

	var updatedData models.Bank
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Model(&bank).Updates(updatedData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bank)
}
