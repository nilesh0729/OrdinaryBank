package api

import "github.com/gin-gonic/gin"

type CreateAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

func (server *Server) CreateAccount(ctx *gin.Context){
	
}