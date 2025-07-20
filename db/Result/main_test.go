package Anuskh

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var TestDb *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:SituBen@localhost:5433/Hiten?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	TestDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to Database:", err)
	}

	testQueries = New(TestDb)

	os.Exit(m.Run())

}
