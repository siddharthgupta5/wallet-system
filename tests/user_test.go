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

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Migrated the schema
	db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})

	return db
}

func TestUserService_CreateUser(t *testing.T) {
	db := setupTestDB()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	t.Run("Create valid user", func(t *testing.T) {
		user, err := userService.CreateUser("John Doe", "john@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("Create user with empty name", func(t *testing.T) {
		_, err := userService.CreateUser("", "john@example.com")
		assert.Error(t, err)
		assert.Equal(t, "name and email are required", err.Error())
	})

	t.Run("Create user with duplicate email", func(t *testing.T) {
		_, err := userService.CreateUser("John Doe", "john@example.com")
		assert.NoError(t, err)

		_, err = userService.CreateUser("Another John", "john@example.com")
		assert.Error(t, err)
		assert.Equal(t, "email already in use", err.Error())
	})
}