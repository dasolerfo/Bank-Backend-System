package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simplebank/db/model"
	token "simplebank/token"

	"github.com/gin-gonic/gin"
)

type argTransferParamsAPI struct {
	FromAccountID int64       `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64       `json:"to_account_id" binding:"required,min=1"`
	Amount        int64       `json:"amount" binding:"required,min=1"`
	Currency      db.Currency `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransferMoneyAPI(ctx *gin.Context) {
	var req argTransferParamsAPI

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	/*fromAccount, err := s.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	toAccount, err := s.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}*/

	fromAccount, val := s.validateAccounts(ctx, req.Currency, req.FromAccountID)
	if !val {
		return
	}

	authPayload := ctx.MustGet(authorizationKey).(*token.Payload)
	owner, err := s.store.GetOwnerByEmail(ctx, authPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if fromAccount.OwnerID != owner.ID {
		ctx.JSON(http.StatusForbidden, errorResponse((errors.New("error: you can't transfer money from an account is not yours"))))
		return
	}

	_, val = s.validateAccounts(ctx, req.Currency, req.ToAccountID)
	if !val {
		return
	}

	if fromAccount.Money <= req.Amount {
		ctx.JSON(http.StatusUnprocessableEntity, "Error: No disposes de suficients diners per realitzar la transferencia")
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

func (s *Server) validateAccounts(ctx *gin.Context, currency db.Currency, accountID int64) (account db.Account, val bool) {

	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		val = false
		return
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusUnprocessableEntity, "Error: Les divises no coincideixen")
		val = false
		return
	}

	return account, true

}
