package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
	"github.com/siddharthgupta5/wallet-api/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransactionService(t *testing.T) {
	db := setupTestDB()
	
	// Created users and wallets
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	user1, _ := userService.CreateUser("User 1", "user1@example.com")
	user2, _ := userService.CreateUser("User 2", "user2@example.com")

	walletRepo := repositories.NewWalletRepository(db)
	walletService := services.NewWalletService(walletRepo, userRepo)
	wallet1, _ := walletService.CreateWallet(user1.ID, "USD")
	wallet2, _ := walletService.CreateWallet(user2.ID, "USD")

	// Seeded some balance
	walletRepo.UpdateBalance(wallet1.ID, 100.0)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, walletRepo)

	t.Run("Create valid credit transaction", func(t *testing.T) {
		txn, err := transactionService.CreateTransaction(wallet1.ID, 50.0, "Test credit", "credit")
		assert.NoError(t, err)
		assert.Equal(t, "credit", txn.Type)
		assert.Equal(t, 50.0, txn.Amount)

		// Checking wallet balance
		wallet, _ := walletService.GetWalletByID(wallet1.ID)
		assert.Equal(t, 150.0, wallet.Balance)
	})

	t.Run("Create valid debit transaction", func(t *testing.T) {
		txn, err := transactionService.CreateTransaction(wallet1.ID, 30.0, "Test debit", "debit")
		assert.NoError(t, err)
		assert.Equal(t, "debit", txn.Type)
		assert.Equal(t, 30.0, txn.Amount)

		// Checking wallet balance
		wallet, _ := walletService.GetWalletByID(wallet1.ID)
		assert.Equal(t, 120.0, wallet.Balance)
	})

	t.Run("Create debit with insufficient balance", func(t *testing.T) {
		_, err := transactionService.CreateTransaction(wallet1.ID, 200.0, "Large debit", "debit")
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
	})

	t.Run("Transfer funds between wallets", func(t *testing.T) {
		debitTxn, creditTxn, err := transactionService.TransferFunds(wallet1.ID, wallet2.ID, 20.0, "Test transfer")
		assert.NoError(t, err)
		assert.Equal(t, "debit", debitTxn.Type)
		assert.Equal(t, "credit", creditTxn.Type)
		assert.Equal(t, 20.0, debitTxn.Amount)
		assert.Equal(t, 20.0, creditTxn.Amount)

		// Checking balances
		wallet1, _ := walletService.GetWalletByID(wallet1.ID)
		wallet2, _ := walletService.GetWalletByID(wallet2.ID)
		assert.Equal(t, 100.0, wallet1.Balance) // 120 - 20 = 100
		assert.Equal(t, 20.0, wallet2.Balance)  // 0 + 20 = 20
	})

	t.Run("Transfer with insufficient balance", func(t *testing.T) {
		_, _, err := transactionService.TransferFunds(wallet1.ID, wallet2.ID, 200.0, "Large transfer")
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance in source wallet", err.Error())
	})
}