package api

import (
	"net/http"
	db "simplebank/db/model"

	"github.com/gin-gonic/gin"
)

// BUSCAR A GO-PLAYGROUND/VALIDATOR

type createAccountRequest struct {
	OwnerID     int64       `json:"owner_id" binding:"required"`
	Currency    db.Currency `json:"currency" binding:"required,oneof=USD EUR JPY KRW"`
	CountryCode int32       `json:"country_code" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		OwnerID:     req.OwnerID,
		Currency:    req.Currency,
		CountryCode: req.CountryCode,
		Money:       0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}
