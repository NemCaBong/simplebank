package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/db/util"
)

var testQueries *Queries

// main entry point for all testing
// inside 1 package -> in this are db
var testDB *sql.DB

func TestMain(m *testing.M) {
	//
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the DB:", err)
	}
	// create a new &Queries struct
	testQueries = New(testDB)

	// start running the unit test
	// return the exit code whether test pass or fail
	// report the code to the os.Exit()
	os.Exit(m.Run())
}
