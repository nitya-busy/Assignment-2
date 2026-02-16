package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OpenAccountRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	HolderRole string `json:"holder_role"`
}

type UpdateAccountRequest struct {
	Type   string  `json:"type" binding:"required,oneof=deposit withdraw"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func OpenSavingsAccount(c *gin.Context) {
	var req OpenAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer models.Customer
	if err := config.GetDB().First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if req.HolderRole == "" {
		req.HolderRole = "primary_holder"
	}

	tx := config.GetDB().Begin()

	account := models.SavingsAccount{
		Balance: 0,
	}

	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customerAccount := models.CustomerAccount{
		CustomerID: req.CustomerID,
		AccountID:  account.ID,
		HolderRole: req.HolderRole,
	}

	if err := tx.Create(&customerAccount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"account":          account,
		"customer_account": customerAccount,
	})
}

func GetAccount(c *gin.Context) {
	id := c.Param("id")

	var account models.SavingsAccount

	if err := config.GetDB().
		Preload("CustomerAccounts.Customer").
		Preload("Transactions").
		First(&account, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}
func UpdateAccount(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	tx := db.Begin()

	var account models.SavingsAccount
	if err := tx.First(&account, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	if req.Type == "deposit" {
		account.Balance += req.Amount
	} else {
		if account.Balance < req.Amount {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}
		account.Balance -= req.Amount
	}

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transaction{
		AccountID: account.ID,
		Type:      req.Type,
		Amount:    req.Amount,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"balance": account.Balance,
	})
}
