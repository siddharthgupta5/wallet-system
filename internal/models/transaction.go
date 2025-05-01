package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	WalletID    uint    `json:"wallet_id" gorm:"not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
	Description string  `json:"description"`
	Type        string  `json:"type" gorm:"not null"` // credit or debit
	Reference   string  `json:"reference" gorm:"unique"`
}