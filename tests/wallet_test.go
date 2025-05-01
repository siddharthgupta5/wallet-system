package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
	"github.com/siddharthgupta5/wallet-api/internal/services"
)

func TestWalletService(t *testing.T) {
	db := setupTestDB()
	
	// Created a user first
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	user, err := userService.CreateUser("Test User", "test@example.com")
	assert.NoError(t, err)

	walletRepo := repositories.NewWalletRepository(db)
	walletService := services.NewWalletService(walletRepo, userRepo)

	t.Run("Create wallet for valid user", func(t *testing.T) {
		wallet, err := walletService.CreateWallet(user.ID, "USD")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, wallet.UserID)
		assert.Equal(t, 0.0, wallet.Balance)
		assert.Equal(t, "USD", wallet.Currency)
	})

	t.Run("Create wallet for non-existent user", func(t *testing.T) {
		_, err := walletService.CreateWallet(999, "USD")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("Get wallet balance", func(t *testing.T) {
		wallet, err := walletService.CreateWallet(user.ID, "USD")
		assert.NoError(t, err)

		balance, err := walletService.GetWalletBalance(wallet.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, balance)
	})

	t.Run("Get non-existent wallet", func(t *testing.T) {
		_, err := walletService.GetWalletByID(999)
		assert.Error(t, err)
	})
}

