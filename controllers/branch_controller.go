package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBranch(c *gin.Context) {
	var branch models.Branch

	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var bank models.Bank
	if err := config.GetDB().First(&bank, branch.BankID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}

	if err := config.GetDB().Create(&branch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, branch)
}
func GetBranch(c *gin.Context) {
	id := c.Param("id")

	var branch models.Branch

	if err := config.GetDB().
		Preload("Bank").
		Preload("Customers").
		First(&branch, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	c.JSON(http.StatusOK, branch)
}
func UpdateBranch(c *gin.Context) {
	id := c.Param("id")

	db := config.GetDB()
	var branch models.Branch
	if err := db.First(&branch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	var updatedData models.Branch
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedData.BankID != 0 {
		var bank models.Bank
		if err := db.First(&bank, updatedData.BankID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
			return
		}
	}

	if err := db.Model(&branch).Updates(updatedData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branch)
}
