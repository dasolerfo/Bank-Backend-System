package api

import (
	"database/sql"
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
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:id binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}
