package main

import (
	"context"
	"database/sql"
	"github.com/mearaj/bhagad-house-booking/common"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"log"
)

var config common.Config
var conn *sql.DB
var store sqlc.Store
var user sqlc.User

func main() {
	initialize()
	createUser()
}

func initialize() {
	config = common.LoadConfig()
	var err error
	conn, err = sql.Open(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	store = sqlc.NewStore(conn)
}

func createUser() {
	hashedPassword, err := utils.HashPassword("12345678")
	if err != nil {
		log.Fatalln(err)
		return
	}
	user, err = store.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:     "Owner Admin",
		Email:    "admin@bhagadhouse.com",
		Password: hashedPassword,
	})
	if err != nil {
		log.Fatalln(err)
	}
	alog.Logger().Println("Successfully created user")
	alog.Logger().Println(user)
}
