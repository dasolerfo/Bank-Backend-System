package api

import (
	"database/sql"
	"log"
	"net/http"
	db "simplebank/db/model"
	"simplebank/factory"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createOwnerRequest struct {
	FirstName     string `json:"first_name" binding:"required"`
	FirstSurname  string `json:"first_surname" binding:"required"`
	SecondSurname string `json:"second_surname" binding:"required"`
	Nationality   int32  `json:"nationality" binding:"required"`
	Password      string `json:"password" binding:"required,min=6"`
	BornAt        string `json:"born_at" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
}

type ownerResponse struct {
	FirstName     string    `json:"first_name"`
	FirstSurname  string    `json:"first_surname"`
	SecondSurname string    `json:"second_surname"`
	Nationality   int32     `json:"nationality"`
	Email         string    `json:"email"`
	BornAt        time.Time `json:"born_at"`
}

func newUserResponse(owner db.Owner) ownerResponse {
	return ownerResponse{
		FirstName:     owner.FirstName,
		FirstSurname:  owner.FirstSurname,
		SecondSurname: owner.SecondSurname,
		Nationality:   owner.Nationality,
		Email:         owner.Email,
		BornAt:        owner.BornAt,
	}
}

func (server *Server) createOwner(ctx *gin.Context) {
	var req createOwnerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	date, err := time.Parse("2006-01-02", req.BornAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	log.Println(date)

	hash_pass, err := factory.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	arg := db.
		CreateOwnerParams{
		FirstName:      req.FirstName,
		FirstSurname:   req.FirstSurname,
		SecondSurname:  req.SecondSurname,
		HashedPassword: hash_pass,
		Nationality:    req.Nationality,
		BornAt:         date,
		Email:          req.Email,
	}

	owner, err := server.store.CreateOwner(ctx, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ownerResponse := newUserResponse(owner)
	ctx.JSON(http.StatusOK, ownerResponse)
}

type loginOwnerRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type loginOwnerResponse struct {
	AccessToken string        `json:"access_token"`
	Owner       ownerResponse `json:"owner"`
}

func (server *Server) loginOwner(ctx *gin.Context) {
	var req loginOwnerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	owner, err := server.store.GetOwnerByEmail(ctx, req.Email)
	if err != nil {

		if err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = factory.CheckPassword(req.Password, owner.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(owner.Email, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginOwnerResponse{
		AccessToken: accessToken,
		Owner:       newUserResponse(owner),
	}

	ctx.JSON(http.StatusOK, response)
}
