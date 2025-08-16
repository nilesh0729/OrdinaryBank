package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
)

type TransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"` //gt == greater than(used in case the amount would be less than 1 but still greater than 0, like Rs0.45)
	Currency      string `json:"currency" binding:"required,oneof=INR USD CAD EUR YEN"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req TransferRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.AccountValidator(ctx,req.FromAccountId,req.Currency){
		return
	}
	if !server.AccountValidator(ctx,req.ToAccountId,req.Currency){
		return
	}

	arg := Anuskh.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID: req.ToAccountId,
		Amount: req.Amount,
	}


	account, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (Server *Server) AccountValidator(ctx *gin.Context, accountID int64, currency string) bool{
	account, err := Server.store.GetAccounts(context.Background(), accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account %d's currency is mismatched : %s vs %s", accountID, currency, account.Currency)	
		ctx.JSON(http.StatusBadRequest, errorResponse(err))	
		return false
	}
	return true
}
