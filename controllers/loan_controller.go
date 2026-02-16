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

type UpdateLoanRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
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
	totalPayableAmount := req.PrincipalAmount + (req.PrincipalAmount * interestRate / 100)

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
		Preload("LoanPayments").
		First(&loan, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	c.JSON(http.StatusOK, loan)
}
func UpdateLoan(c *gin.Context) {
	id := c.Param("id")

	var req UpdateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	tx := db.Begin()

	var loan models.Loan
	if err := tx.First(&loan, id).Error; err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount exceeds pending balance"})
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
		"message":        "Loan updated successfully",
		"pending_amount": loan.PendingAmount,
		"status":         loan.Status,
	})
}
