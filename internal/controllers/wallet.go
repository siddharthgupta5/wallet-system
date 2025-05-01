package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wallet-api/internal/models"
	"github.com/yourusername/wallet-api/internal/services"
	"github.com/yourusername/wallet-api/internal/utils"
)

type WalletController struct {
	walletService services.WalletService
}

func NewWalletController(walletService services.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

type CreateWalletRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Currency string `json:"currency"`
}

func (c *WalletController) CreateWallet(ctx *gin.Context) {
	var req CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	wallet, err := c.walletService.CreateWallet(req.UserID, req.Currency)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, wallet)
}

func (c *WalletController) GetWallet(ctx *gin.Context) {
	walletID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := c.walletService.GetWalletByID(walletID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Wallet not found")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, wallet)
}

func (c *WalletController) GetWalletBalance(ctx *gin.Context) {
	walletID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	balance, err := c.walletService.GetWalletBalance(walletID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Wallet not found")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"balance": balance})
}

func (c *WalletController) GetUserWallets(ctx *gin.Context) {
	userID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	wallets, err := c.walletService.GetUserWallets(userID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "User not found or no wallets")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, wallets)
}