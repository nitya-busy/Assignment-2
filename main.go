package main

import (
	"banking-system/config"
	"banking-system/models"
	"banking-system/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database configuration
	dbConfig := config.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Set default values if env variables are not set
	if dbConfig.Host == "" {
		dbConfig.Host = "localhost"
	}
	if dbConfig.Port == "" {
		dbConfig.Port = "5432"
	}
	if dbConfig.User == "" {
		dbConfig.User = "postgres"
	}
	if dbConfig.Password == "" {
		dbConfig.Password = "postgres"
	}
	if dbConfig.DBName == "" {
		dbConfig.DBName = "banking_system"
	}
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	// Initialize database
	if err := config.InitDB(dbConfig); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	// Run migrations
	db := config.GetDB()
	
	// Drop all tables to ensure clean migration (for development)
	db.Migrator().DropTable(
		&models.LoanPayment{},
		&models.Loan{},
		&models.Transaction{},
		&models.CustomerAccount{},
		&models.SavingsAccount{},
		&models.Customer{},
		&models.Branch{},
		&models.Bank{},
	)
	
	if err := db.AutoMigrate(
		&models.Bank{},
		&models.Branch{},
		&models.Customer{},
		&models.SavingsAccount{},
		&models.CustomerAccount{},
		&models.Transaction{},
		&models.Loan{},
		&models.LoanPayment{},
	); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	log.Println("Database migrations completed successfully")

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Add a health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Banking System API is running",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Banking System API on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
