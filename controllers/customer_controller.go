package controllers

import (
	"banking-system/config"
	"banking-system/models"
	"fmt"
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

	if result := config.GetDB().Create(&customer); result.Error != nil {
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

func ListCustomersWithAccounts(c *gin.Context) {

	fmt.Println("DEBUG: ListCustomersWithAccounts API HIT")

	type CustomerAccountDetail struct {
		AccountNo  uint    `json:"account_no"`
		BranchName string  `json:"branch_name"`
		BankName   string  `json:"bank_name"`
		Balance    float64 `json:"balance"`
	}

	type CustomerWithAccountsResponse struct {
		ID       uint                    `json:"id"`
		Name     string                  `json:"name"`
		Email    string                  `json:"email"`
		Phone    string                  `json:"phone"`
		Accounts []CustomerAccountDetail `json:"accounts"`
	}

	query := `
		SELECT 
			c.id as customer_id,
			c.name as customer_name,
			c.email as customer_email,
			c.phone as customer_phone,
			sa.id as account_no,
			br.name as branch_name,
			b.name as bank_name,
			sa.balance as balance
		FROM customers c
		LEFT JOIN customer_accounts ca ON c.id = ca.customer_id
		LEFT JOIN savings_accounts sa ON ca.account_id = sa.id
		LEFT JOIN branches br ON c.branch_id = br.id
		LEFT JOIN banks b ON br.bank_id = b.id
		ORDER BY c.id, sa.id
	`

	type CustomerQueryResult struct {
		CustomerID    uint
		CustomerName  string
		CustomerEmail string
		CustomerPhone string
		AccountNo     *uint
		BranchName    *string
		BankName      *string
		Balance       *float64
	}

	var results []CustomerQueryResult

	if err := config.GetDB().Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customerMap := make(map[uint]*CustomerWithAccountsResponse)

	for _, row := range results {

		if _, exists := customerMap[row.CustomerID]; !exists {
			customerMap[row.CustomerID] = &CustomerWithAccountsResponse{
				ID:       row.CustomerID,
				Name:     row.CustomerName,
				Email:    row.CustomerEmail,
				Phone:    row.CustomerPhone,
				Accounts: []CustomerAccountDetail{},
			}
		}

		if row.AccountNo != nil && row.BranchName != nil && row.BankName != nil && row.Balance != nil {
			customerMap[row.CustomerID].Accounts = append(
				customerMap[row.CustomerID].Accounts,
				CustomerAccountDetail{
					AccountNo:  *row.AccountNo,
					BranchName: *row.BranchName,
					BankName:   *row.BankName,
					Balance:    *row.Balance,
				},
			)
		}
	}

	response := make([]CustomerWithAccountsResponse, 0, len(customerMap))
	for _, customer := range customerMap {
		response = append(response, *customer)
	}

	c.JSON(http.StatusOK, response)
}
