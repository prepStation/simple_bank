package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/prepStation/simple_bank/utils"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database %v\n", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
