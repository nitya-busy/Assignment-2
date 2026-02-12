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
	if result := config.GetDB().First(&branch, customer.BranchID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	result := config.GetDB().Create(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}
func GetCustomer(c *gin.Context) {
	id := c.Param("customer_id")
	var customer models.Customer

	result := config.GetDB().
		Preload("Branch").
		Preload("CustomerAccounts.Account").
		Preload("Loans").
		First(&customer, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
func GetCustomersByBranch(c *gin.Context) {
	branchID := c.Param("branch_id")
	var customers []models.Customer

	result := config.GetDB().
		Where("branch_id = ?", branchID).
		Preload("CustomerAccounts.Account").
		Preload("Loans").
		Find(&customers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}
