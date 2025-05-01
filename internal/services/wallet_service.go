package services

import (
	"errors"
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
)

type WalletService interface {
	CreateWallet(userID uint, currency string) (*models.Wallet, error)
	GetWalletByID(walletID uint) (*models.Wallet, error)
	GetUserWallets(userID uint) ([]models.Wallet, error)
	GetWalletBalance(walletID uint) (float64, error)
}

type walletService struct {
	walletRepo repositories.WalletRepository
	userRepo   repositories.UserRepository
}

func NewWalletService(walletRepo repositories.WalletRepository, userRepo repositories.UserRepository) WalletService {
	return &walletService{
		walletRepo: walletRepo,
		userRepo:   userRepo,
	}
}

func (s *walletService) CreateWallet(userID uint, currency string) (*models.Wallet, error) {
	if currency == "" {
		currency = "USD"
	}

	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	wallet := &models.Wallet{
		UserID:   userID,
		Currency: currency,
		Balance:  0,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetWalletByID(walletID uint) (*models.Wallet, error) {
	return s.walletRepo.FindByID(walletID)
}

func (s *walletService) GetUserWallets(userID uint) ([]models.Wallet, error) {
	return s.walletRepo.FindByUserID(userID)
}

func (s *walletService) GetWalletBalance(walletID uint) (float64, error) {
	wallet, err := s.walletRepo.FindByID(walletID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}