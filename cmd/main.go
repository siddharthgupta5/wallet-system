package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/siddharthgupta5/wallet-api/internal/config"
	"github.com/siddharthgupta5/wallet-api/internal/controllers"
	"github.com/siddharthgupta5/wallet-api/internal/db"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
	"github.com/siddharthgupta5/wallet-api/internal/routes"
	"github.com/siddharthgupta5/wallet-api/internal/services"
)

func main() {
	// Loading configuration
	cfg := config.LoadConfig()

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Initialized repositories
	userRepo := repositories.NewUserRepository(database)
	walletRepo := repositories.NewWalletRepository(database)
	transactionRepo := repositories.NewTransactionRepository(database)

	// Initialized services
	userService := services.NewUserService(userRepo)
	walletService := services.NewWalletService(walletRepo, userRepo)
	transactionService := services.NewTransactionService(transactionRepo, walletRepo)

	// Initialized controllers
	userController := controllers.NewUserController(userService)
	walletController := controllers.NewWalletController(walletService)
	transactionController := controllers.NewTransactionController(transactionService)

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router, userController, walletController, transactionController)

	// Starting server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}