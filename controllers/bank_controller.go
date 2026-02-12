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

	result := config.GetDB().Create(&bank)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, bank)
}
func GetBank(c *gin.Context) {
	id := c.Param("bank_id")
	var bank models.Bank

	result := config.GetDB().Preload("Branches").First(&bank, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}

	c.JSON(http.StatusOK, bank)
}
func GetAllBanks(c *gin.Context) {
	var banks []models.Bank

	result := config.GetDB().Preload("Branches").Find(&banks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, banks)
}
