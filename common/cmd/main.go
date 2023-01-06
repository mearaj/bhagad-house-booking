package main

import (
	"context"
	"database/sql"
	"github.com/mearaj/bhagad-house-booking/backend"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"log"
)

func main() {
	config := backend.LoadConfig()
	conn, err := sql.Open(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		log.Fatalln("cannot connect to db:", err)
	}
	store := sqlc.NewStore(conn)
	hashedPassword, err := utils.HashPassword("12345678")
	if err != nil {
		log.Fatalln(err)
		return
	}
	user, err := store.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:     "Owner Admin",
		Email:    "admin@bhagadhouse.com",
		Password: hashedPassword,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("User created")
	log.Println(user)
}
