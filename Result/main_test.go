package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:SituBen@localhost:5433/Hiten?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to Database:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())

}
