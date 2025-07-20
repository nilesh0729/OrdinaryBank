package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
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
