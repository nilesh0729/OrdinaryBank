package api

import (
	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
)

type Server struct {
	store  *Anuskh.Store
	router *gin.Engine
}

func NewServer(store *Anuskh.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.store.CreateAccount)

	server.router = router

	return server
}
