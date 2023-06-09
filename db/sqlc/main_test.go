package db

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
	dbSource = "postgresql://root:thobeogalaxy257@localhost:5432/simple_bank?sslmode=disable"
)

// main entry point for all testing
// inside 1 package -> in this are db
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
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
