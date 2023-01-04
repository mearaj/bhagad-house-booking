package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/backend"
	"github.com/mearaj/bhagad-house-booking/backend/api"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := backend.LoadConfig()
	conn, err := sql.Open(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln(err)
	}
	err = server.Start()
	if err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
