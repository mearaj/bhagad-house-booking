package sqlc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config := LoadConfig()
	var err error
	testDB, err = sql.Open(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
