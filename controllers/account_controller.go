package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OpenAccountRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	HolderRole string `json:"holder_role"` // primary_holder or joint_holder, defaults to primary_holder
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

// OpenSavingsAccount opens a new savings account for a customer
func OpenSavingsAccount(c *gin.Context) {
	var req OpenAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify customer exists
	var customer models.Customer
	if result := config.GetDB().First(&customer, req.CustomerID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Set default holder role
	if req.HolderRole == "" {
		req.HolderRole = "primary_holder"
	}

	tx := config.GetDB().Begin()

	// Create new savings account
	account := models.SavingsAccount{
		Balance: 0,
	}

	if result := tx.Create(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Link customer to account via CustomerAccount
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

// AddAccountHolder adds a customer as a joint holder to an existing account
func AddAccountHolder(c *gin.Context) {
	accountID := c.Param("account_id")
	var req AddAccountHolderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify account exists
	var account models.SavingsAccount
	if result := config.GetDB().First(&account, accountID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Verify customer exists
	var customer models.Customer
	if result := config.GetDB().First(&customer, req.CustomerID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Check if customer is already linked to this account
	var existingLink models.CustomerAccount
	if result := config.GetDB().Where("customer_id = ? AND account_id = ?", req.CustomerID, account.ID).First(&existingLink); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer is already linked to this account"})
		return
	}

	// Link customer to account
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

// GetAccount retrieves account details with all holders
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

// Deposit deposits money into an account
func Deposit(c *gin.Context) {
	id := c.Param("account_id")
	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.GetDB().Begin()

	// Get account
	var account models.SavingsAccount
	if result := tx.First(&account, id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Update balance
	account.Balance += req.Amount
	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Record transaction
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

// Withdraw withdraws money from an account
func Withdraw(c *gin.Context) {
	id := c.Param("account_id")
	var req WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.GetDB().Begin()

	// Get account
	var account models.SavingsAccount
	if result := tx.First(&account, id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Check balance
	if account.Balance < req.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Update balance
	account.Balance -= req.Amount
	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Record transaction
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

// GetTransactions retrieves all transactions for an account
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
