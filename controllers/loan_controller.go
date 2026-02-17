package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TakeLoanRequest struct {
	CustomerID      uint    `json:"customer_id" binding:"required"`
	LoanType        string  `json:"loan_type" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
}

type RepayLoanRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
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
	totalPayableAmount := req.PrincipalAmount + (req.PrincipalAmount * interestRate / 100.0)

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
		Preload("Customer").
		Preload("Customer.Branch").
		Preload("Customer.Branch.Bank").
		Preload("LoanPayments").
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

	var loan models.Loan

	if err := tx.
		Preload("Customer").
		Preload("Customer.Branch").
		Preload("Customer.Branch.Bank").
		Preload("LoanPayments").
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

	loan.PendingAmount -= req.Amount

	if loan.PendingAmount == 0 {
		loan.Status = "CLOSED"
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

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Loan repayment successful",
		"loan":    loan,
		"payment": payment,
	})
}

func GetLoanInterest(c *gin.Context) {
	loanID := c.Param("id")
	var loan models.Loan

	if err := config.GetDB().First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	interestForThisYear := (loan.PendingAmount * loan.InterestRate) / 100.0

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
		Preload("Loan").
		Preload("Loan.Customer").
		Preload("Loan.Customer.Branch").
		Preload("Loan.Customer.Branch.Bank").
		Order("payment_date DESC").
		Find(&payments).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
