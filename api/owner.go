package api

import (
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

type createOwnerResponse struct {
	FirstName     string    `json:"first_name"`
	FirstSurname  string    `json:"first_surname"`
	SecondSurname string    `json:"second_surname"`
	Nationality   int32     `json:"nationality"`
	Email         string    `json:"email"`
	BornAt        time.Time `json:"born_at"`
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

	ownerResponse := createOwnerResponse{
		FirstName:     owner.FirstName,
		FirstSurname:  owner.FirstSurname,
		SecondSurname: owner.SecondSurname,
		Nationality:   owner.Nationality,
		Email:         owner.Email,
		BornAt:        owner.BornAt,
	}

	ctx.JSON(http.StatusOK, ownerResponse)
}
