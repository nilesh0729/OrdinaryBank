package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nilesh0729/OrdinaryBank/api"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
	"github.com/nilesh0729/OrdinaryBank/util"
)


func main(){
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load Config", err)
	}
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot Connect To db: ", err)
	}


	store := Anuskh.NewTxConn(conn)
	server := api.NewServer(store)


	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("Cannot Start Server : ", err)
	}
}