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
	dbConfig := config.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
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
	if err := config.InitDB(dbConfig); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	db := config.GetDB()
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
	router := gin.Default()
	routes.SetupRoutes(router)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Banking System API is running",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Banking System API on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
