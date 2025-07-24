package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
)

type CreateAccountRequest struct {
	Owner string `json:"owner" binding:"required"`

	Currency string `json:"currency" binding:"required,oneof=INR USD CAD EUR YEN"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := Anuskh.CreateAccountsParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}
//
//
type GetAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccounts(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}


	ctx.JSON(http.StatusOK, account)
}
//
//
//
type ListAccountRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}


func (server *Server) ListAccount(ctx *gin.Context) {
	var req ListAccountRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	arg := Anuskh.ListAccountsParams{
		Limit: (req.PageSize),
		Offset: (req.PageID - 1)* req.PageSize ,
	}

	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
