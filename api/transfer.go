package api

import (
	"net/http"
	db "simplebank/db/model"

	"github.com/gin-gonic/gin"
)

type argTransferParamsAPI struct {
	FromAccountID int64       `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64       `json:"to_account_id" binding:"required,min=1"`
	Amount        int64       `json:"amount" binding:"required,min=1"`
	Currency      db.Currency `json:"currency" binding:"required,oneof=USD EUR JPY KRW"`
}

func (s *Server) createTransferMoneyAPI(ctx *gin.Context) {
	var req argTransferParamsAPI

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}
