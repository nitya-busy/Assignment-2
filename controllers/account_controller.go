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

type AddAccountHolderRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	HolderRole string `json:"holder_role" binding:"required"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func OpenSavingsAccount(c *gin.Context) {
	var req OpenAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var customer models.Customer
	if result := config.GetDB().First(&customer, req.CustomerID); result.Error != nil {
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
	if result := tx.Create(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	customerAccount := models.CustomerAccount{
		CustomerID: req.CustomerID,
		AccountID:  account.ID,
		HolderRole: req.HolderRole,
	}
	if result := tx.Create(&customerAccount); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"account":          account,
		"customer_account": customerAccount,
	})
}
func AddAccountHolder(c *gin.Context) {
	accountID := c.Param("account_id")
	var req AddAccountHolderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var account models.SavingsAccount
	if result := config.GetDB().First(&account, accountID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	var customer models.Customer
	if result := config.GetDB().First(&customer, req.CustomerID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	var existingLink models.CustomerAccount
	if result := config.GetDB().Where("customer_id = ? AND account_id = ?", req.CustomerID, account.ID).First(&existingLink); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer is already linked to this account"})
		return
	}
	customerAccount := models.CustomerAccount{
		CustomerID: req.CustomerID,
		AccountID:  account.ID,
		HolderRole: req.HolderRole,
	}

	result := config.GetDB().Create(&customerAccount)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, customerAccount)
}
func GetAccount(c *gin.Context) {
	id := c.Param("account_id")
	var account models.SavingsAccount

	result := config.GetDB().
		Preload("CustomerAccounts.Customer").
		Preload("Transactions").
		First(&account, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}
func Deposit(c *gin.Context) {
	id := c.Param("account_id")
	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.GetDB().Begin()
	var account models.SavingsAccount
	if result := tx.First(&account, id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	account.Balance += req.Amount
	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	transaction := models.Transaction{
		AccountID: account.ID,
		Type:      "DEPOSIT",
		Amount:    req.Amount,
	}

	if result := tx.Create(&transaction); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"account":     account,
		"transaction": transaction,
	})
}
func Withdraw(c *gin.Context) {
	id := c.Param("account_id")
	var req WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.GetDB().Begin()
	var account models.SavingsAccount
	if result := tx.First(&account, id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	if account.Balance < req.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}
	account.Balance -= req.Amount
	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	transaction := models.Transaction{
		AccountID: account.ID,
		Type:      "WITHDRAW",
		Amount:    req.Amount,
	}

	if result := tx.Create(&transaction); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"account":     account,
		"transaction": transaction,
	})
}
func GetTransactions(c *gin.Context) {
	id := c.Param("account_id")
	var transactions []models.Transaction

	result := config.GetDB().Where("account_id = ?", id).Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
