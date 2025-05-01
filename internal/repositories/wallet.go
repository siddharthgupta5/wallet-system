package repositories

import (
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByID(id uint) (*models.Wallet, error)
	FindByUserID(userID uint) ([]models.Wallet, error)
	UpdateBalance(walletID uint, amount float64) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, id).Error
	return &wallet, err
}

func (r *walletRepository) FindByUserID(userID uint) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *walletRepository) UpdateBalance(walletID uint, amount float64) error {
	return r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}