package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type TakeLoanRequest struct {
	CustomerID      uint    `json:"customer_id" binding:"required"`
	LoanType        string  `json:"loan_type" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
}

type RepayLoanRequest struct {
	AccountID uint    `json:"account_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

type InterestResponse struct {
	LoanID              uint      `json:"loan_id"`
	YearlyInterestRate  float64   `json:"yearly_interest_rate"`
	PendingAmount       float64   `json:"pending_amount"`
	InterestForThisYear float64   `json:"interest_for_this_year"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

func TakeLoan(c *gin.Context) {
	var req TakeLoanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer models.Customer
	if err := config.GetDB().First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	interestRate := 12.0
	totalPayableAmount := req.PrincipalAmount +
		(req.PrincipalAmount * interestRate / 100.0)

	loan := models.Loan{
		CustomerID:         req.CustomerID,
		LoanType:           req.LoanType,
		PrincipalAmount:    req.PrincipalAmount,
		InterestRate:       interestRate,
		TotalPayableAmount: totalPayableAmount,
		PendingAmount:      totalPayableAmount,
		StartDate:          time.Now(),
		Status:             "ACTIVE",
	}

	if err := config.GetDB().Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func GetLoan(c *gin.Context) {
	id := c.Param("id")
	var loan models.Loan

	if err := config.GetDB().
		Preload("Customer").
		Preload("Customer.Branch").
		Preload("Customer.Branch.Bank").
		Preload("LoanPayments").
		First(&loan, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	c.JSON(http.StatusOK, loan)
}

func GetCustomerLoans(c *gin.Context) {
	customerID := c.Param("customer_id")
	var loans []models.Loan

	if err := config.GetDB().
		Where("customer_id = ?", customerID).
		Preload("LoanPayments").
		Order("start_date DESC").
		Find(&loans).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}

func RepayLoan(c *gin.Context) {
	loanID := c.Param("id")

	var req RepayLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.GetDB().Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var loan models.Loan
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&loan, loanID).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	if loan.Status == "CLOSED" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan already closed"})
		return
	}

	if req.Amount > loan.PendingAmount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount exceeds pending amount"})
		return
	}

	var account models.SavingsAccount
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&account, req.AccountID).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var customerAccount models.CustomerAccount
	if err := tx.
		Where("customer_id = ? AND account_id = ?", loan.CustomerID, account.ID).
		First(&customerAccount).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account does not belong to loan customer"})
		return
	}

	if account.Balance < req.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	account.Balance -= req.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loan.PendingAmount -= req.Amount
	if loan.PendingAmount == 0 {
		loan.Status = "CLOSED"
		now := time.Now()
		loan.EndDate = &now
	}

	if err := tx.Save(&loan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payment := models.LoanPayment{
		LoanID:      loan.ID,
		Amount:      req.Amount,
		PaymentDate: time.Now(),
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transaction{
		AccountID: account.ID,
		Type:      "WITHDRAW",
		Amount:    req.Amount,
		Balance:   account.Balance,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"withdraw_amount": req.Amount,
		"updated_balance": account.Balance,
		"transaction":     transaction,
	})
}

func GetLoanInterest(c *gin.Context) {
	loanID := c.Param("id")
	var loan models.Loan

	if err := config.GetDB().First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	interestForThisYear :=
		(loan.PendingAmount * loan.InterestRate) / 100.0

	response := InterestResponse{
		LoanID:              loan.ID,
		YearlyInterestRate:  loan.InterestRate,
		PendingAmount:       loan.PendingAmount,
		InterestForThisYear: interestForThisYear,
		CalculatedAt:        time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

func GetLoanPayments(c *gin.Context) {
	loanID := c.Param("id")
	var payments []models.LoanPayment

	if err := config.GetDB().
		Where("loan_id = ?", loanID).
		Order("payment_date DESC").
		Find(&payments).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
