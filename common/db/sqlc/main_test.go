package sqlc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	_, p, _, ok := runtime.Caller(0) // provides path of this main file
	if !ok {
		log.Fatalln("error in runtime.Caller, cannot load path")
	}
	p = filepath.Join(p, filepath.FromSlash("../../.."))
	config, err := utils.LoadConfig(p)
	if err != nil {
		log.Fatalln("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
