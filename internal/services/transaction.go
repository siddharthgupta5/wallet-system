package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(walletID uint, amount float64, description, txnType string) (*models.Transaction, error)
	TransferFunds(sourceWalletID, destinationWalletID uint, amount float64, description string) (*models.Transaction, *models.Transaction, error)
	GetWalletTransactions(walletID uint) ([]models.Transaction, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	walletRepo     repositories.WalletRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, walletRepo repositories.WalletRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		walletRepo:     walletRepo,
	}
}

func (s *transactionService) CreateTransaction(walletID uint, amount float64, description, txnType string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if txnType != "credit" && txnType != "debit" {
		return nil, errors.New("invalid transaction type")
	}

	wallet, err := s.walletRepo.FindByID(walletID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// For debit transactions, checking sufficient balance
	if txnType == "debit" && wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Update wallet balance
	balanceChange := amount
	if txnType == "debit" {
		balanceChange = -amount
	}

	if err := s.walletRepo.UpdateBalance(walletID, balanceChange); err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		WalletID:    walletID,
		Amount:      amount,
		Description: description,
		Type:        txnType,
		Reference:   uuid.New().String(),
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		// Revert balance update if transaction creation fails
		_ = s.walletRepo.UpdateBalance(walletID, -balanceChange)
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) TransferFunds(sourceWalletID, destinationWalletID uint, amount float64, description string) (*models.Transaction, *models.Transaction, error) {
	if sourceWalletID == destinationWalletID {
		return nil, nil, errors.New("cannot transfer to the same wallet")
	}

	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	// Round to 2 decimal places to avoid floating point precision issues
	amount = math.Round(amount*100) / 100

	// Get source wallet
	sourceWallet, err := s.walletRepo.FindByID(sourceWalletID)
	if err != nil {
		return nil, nil, errors.New("source wallet not found")
	}

	// Checking sufficient balance
	if sourceWallet.Balance < amount {
		return nil, nil, errors.New("insufficient balance in source wallet")
	}

	// Get destination wallet
	_, err = s.walletRepo.FindByID(destinationWalletID)
	if err != nil {
		return nil, nil, errors.New("destination wallet not found")
	}

	// Creating debit transaction for source wallet
	debitDescription := fmt.Sprintf("Transfer to wallet %d: %s", destinationWalletID, description)
	debitTxn, err := s.CreateTransaction(sourceWalletID, amount, debitDescription, "debit")
	if err != nil {
		return nil, nil, err
	}

	// Creating credit transaction for destination wallet
	creditDescription := fmt.Sprintf("Transfer from wallet %d: %s", sourceWalletID, description)
	creditTxn, err := s.CreateTransaction(destinationWalletID, amount, creditDescription, "credit")
	if err != nil {
		// If credit fails, revert the debit
		_ = s.walletRepo.UpdateBalance(sourceWalletID, amount)
		_ = s.transactionRepo.Create(&models.Transaction{
			WalletID:    sourceWalletID,
			Amount:      amount,
			Description: "Reverted transfer due to failure",
			Type:        "credit",
			Reference:   uuid.New().String(),
		})
		return nil, nil, err
	}

	return debitTxn, creditTxn, nil
}

func (s *transactionService) GetWalletTransactions(walletID uint) ([]models.Transaction, error) {
	return s.transactionRepo.FindByWalletID(walletID)
}