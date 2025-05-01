package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/siddharthgupta5/wallet-api/internal/controllers"
)

func SetupRoutes(router *gin.Engine, userController *controllers.UserController, walletController *controllers.WalletController, transactionController *controllers.TransactionController) {
	api := router.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
		}

		// Wallet routes
		wallets := api.Group("/wallets")
		{
			wallets.POST("", walletController.CreateWallet)
			wallets.GET("/:id", walletController.GetWallet)
			wallets.GET("/:id/balance", walletController.GetWalletBalance)
			wallets.GET("/user/:id", walletController.GetUserWallets)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionController.CreateTransaction)
			transactions.POST("/transfer", transactionController.TransferFunds)
			transactions.GET("/wallet/:id", transactionController.GetWalletTransactions)
		}
	}
}