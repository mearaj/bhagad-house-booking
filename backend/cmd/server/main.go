package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/backend/api"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := utils.LoadConfig()
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln(err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
