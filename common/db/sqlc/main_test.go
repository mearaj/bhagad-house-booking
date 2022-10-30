package sqlc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

const defaultDatabaseDriver = "postgres"
const defaultDatabaseURL = "postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable"
const defaultServerAddress = "0.0.0.0:8080"
const defaultTokenSymmetricKey = "12345678901234567890123456789012"
const defaultAccessTokenDuration = "15m"

func TestMain(m *testing.M) {
	config := utils.LoadConfig()
	config.DBSource = defaultDatabaseURL
	config.DBDriver = defaultDatabaseDriver
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
