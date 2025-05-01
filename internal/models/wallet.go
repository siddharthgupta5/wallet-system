package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID      uint         `json:"user_id" gorm:"not null"`
	Balance     float64      `json:"balance" gorm:"not null;default:0"`
	Currency    string       `json:"currency" gorm:"not null;default:'USD'"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:WalletID"`
}