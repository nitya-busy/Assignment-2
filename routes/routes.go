package routes

import (
	"banking-system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/banks", controllers.CreateBank)
	router.GET("/banks/:id", controllers.GetBank)
	router.PUT("/banks/:id", controllers.UpdateBank)

	router.POST("/branches", controllers.CreateBranch)
	router.GET("/branches/:id", controllers.GetBranch)
	router.PUT("/branches/:id", controllers.UpdateBranch)

	router.POST("/customers", controllers.CreateCustomer)
	router.GET("/customers/:id", controllers.GetCustomer)
	router.PUT("/customers/:id", controllers.UpdateCustomer)

	router.POST("/accounts", controllers.OpenSavingsAccount)
	router.GET("/accounts/:id", controllers.GetAccount)
	router.PUT("/accounts/:id", controllers.UpdateAccount)

	router.POST("/loans", controllers.TakeLoan)
	router.GET("/loans/:id", controllers.GetLoan)
	router.PUT("/loans/:id", controllers.UpdateLoan)
}
