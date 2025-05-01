package repositories

import (
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByWalletID(walletID uint) ([]models.Transaction, error)
	FindByReference(reference string) (*models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByWalletID(walletID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("wallet_id = ?", walletID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByReference(reference string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("reference = ?", reference).First(&transaction).Error
	return &transaction, err
}