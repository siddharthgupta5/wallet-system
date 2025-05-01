package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siddharthgupta5/wallet-api/internal/services"
	"github.com/siddharthgupta5/wallet-api/internal/utils"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

type CreateTransactionRequest struct {
	WalletID    uint    `json:"wallet_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	Type        string  `json:"type" binding:"required,oneof=credit debit"`
}

type TransferRequest struct {
	SourceWalletID      uint    `json:"source_wallet_id" binding:"required"`
	DestinationWalletID uint    `json:"destination_wallet_id" binding:"required"`
	Amount              float64 `json:"amount" binding:"required,gt=0"`
	Description         string  `json:"description"`
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	transaction, err := c.transactionService.CreateTransaction(req.WalletID, req.Amount, req.Description, req.Type)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, transaction)
}

func (c *TransactionController) TransferFunds(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	debitTxn, creditTxn, err := c.transactionService.TransferFunds(
		req.SourceWalletID,
		req.DestinationWalletID,
		req.Amount,
		req.Description,
	)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, gin.H{
		"debit_transaction":  debitTxn,
		"credit_transaction": creditTxn,
		"message":            "Transfer completed successfully",
	})
}

func (c *TransactionController) GetWalletTransactions(ctx *gin.Context) {
	walletID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	transactions, err := c.transactionService.GetWalletTransactions(walletID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Wallet not found or no transactions")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, transactions)
}
