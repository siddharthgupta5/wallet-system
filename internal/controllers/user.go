package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"github.com/siddharthgupta5/wallet-api/internal/services"
	"github.com/siddharthgupta5/wallet-api/internal/utils"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.userService.CreateUser(req.Name, req.Email)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, user)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, user)
}