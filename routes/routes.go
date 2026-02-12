package routes

import (
	"banking-system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/banks", controllers.CreateBank)
	router.GET("/banks/:bank_id", controllers.GetBank)
	router.GET("/banks", controllers.GetAllBanks)

	router.POST("/branches", controllers.CreateBranch)
	router.GET("/branches/:branch_id", controllers.GetBranch)
	router.GET("/banks/:bank_id/branches", controllers.GetBranchesByBank)

	router.POST("/customers", controllers.CreateCustomer)
	router.GET("/customers/:customer_id", controllers.GetCustomer)
	router.GET("/branches/:branch_id/customers", controllers.GetCustomersByBranch)

	router.POST("/accounts/savings", controllers.OpenSavingsAccount)
	router.POST("/accounts/:account_id/holders", controllers.AddAccountHolder)
	router.GET("/accounts/:account_id", controllers.GetAccount)
	router.POST("/accounts/:account_id/deposit", controllers.Deposit)
	router.POST("/accounts/:account_id/withdraw", controllers.Withdraw)
	router.GET("/accounts/:account_id/transactions", controllers.GetTransactions)

	router.POST("/loans", controllers.TakeLoan)
	router.GET("/loans/:id", controllers.GetLoan)
	router.GET("/customers/:customer_id/loans", controllers.GetCustomerLoans)
	router.POST("/loans/:id/repay", controllers.RepayLoan)
	router.GET("/loans/:id/interest", controllers.GetLoanInterest)
	router.GET("/loans/:id/payments", controllers.GetLoanPayments)
}
