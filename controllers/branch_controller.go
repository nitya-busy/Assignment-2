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
	if result := config.GetDB().First(&bank, branch.BankID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}

	result := config.GetDB().Create(&branch)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, branch)
}
func GetBranch(c *gin.Context) {
	id := c.Param("branch_id")
	var branch models.Branch

	result := config.GetDB().Preload("Bank").Preload("Customers").First(&branch, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	c.JSON(http.StatusOK, branch)
}
func GetBranchesByBank(c *gin.Context) {
	bankID := c.Param("bank_id")
	var branches []models.Branch

	result := config.GetDB().Where("bank_id = ?", bankID).Preload("Customers").Find(&branches)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, branches)
}
