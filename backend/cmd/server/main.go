package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/backend/api"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
)

func main() {
	_, p, _, ok := runtime.Caller(0) // provides path of this main file
	if !ok {
		log.Fatalln("error in runtime.Caller, cannot load path")
	}
	p = filepath.Join(p, filepath.FromSlash("../../../"))
	config, err := utils.LoadConfig(p)
	if err != nil {
		log.Fatalln("cannot load config:", err)
	}
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
