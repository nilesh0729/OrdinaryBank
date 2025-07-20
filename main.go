package main

import (
	"database/sql"
	"log"

    _"github.com/lib/pq"
	"github.com/nilesh0729/OrdinaryBank/api"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
)

const(
	dbDriver = "postgres"
	dbSource = "postgres://root:SituBen@localhost:5433/Hiten?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main(){
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot Connect To db: ", err)
	}


	store := Anuskh.NewTxConn(conn)
	server := api.NewServer(store)


	err = server.Start(serverAddress)
	if err != nil{
		log.Fatal("Cannot Start Server : ", err)
	}
}