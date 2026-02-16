package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var branch models.Branch
	if err := config.GetDB().First(&branch, customer.BranchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	if err := config.GetDB().Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}
func GetCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer

	if err := config.GetDB().
		Preload("Branch").
		Preload("CustomerAccounts.Account").
		Preload("Loans").
		First(&customer, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	db := config.GetDB()
	var customer models.Customer
	if err := db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var updatedData models.Customer
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedData.BranchID != 0 {
		var branch models.Branch
		if err := db.First(&branch, updatedData.BranchID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
	}

	if err := db.Model(&customer).Updates(updatedData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}
