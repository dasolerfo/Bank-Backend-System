package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simplebank/db/model"
	token "simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// IMPORTANT!!!!

// CERCAR A GO-PLAYGROUND/VALIDATOR

type createAccountRequest struct {
	Currency    db.Currency `json:"currency" binding:"required,currency"`
	CountryCode int32       `json:"country_code" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	owner, err := server.store.GetOwnerByEmail(ctx, authPayload.Email)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		OwnerID:     owner.ID,
		Currency:    req.Currency,
		CountryCode: req.CountryCode,
		Money:       0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "foreing_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	owner, err := server.store.GetOwnerByEmail(ctx, authPayload.Email)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.OwnerID != owner.ID {
		err := errors.New("the account does not belong to the owner")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	owner, err := server.store.GetOwnerByEmail(ctx, authPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListAccountParams{
		OwnerID: owner.ID,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccount(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)

}
