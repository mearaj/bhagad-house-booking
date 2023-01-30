package main

import (
	"context"
	"github.com/mearaj/bhagad-house-booking/backend"
	"github.com/mearaj/bhagad-house-booking/backend/api"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	config := backend.LoadConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DatabaseURL))
	if err != nil {
		log.Fatalln("could not connect to db:", err)
		return
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalln("could not connect to db:", err)
		return
	}
	server, err := api.NewServer(config, client)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = server.Start()
	if err != nil {
		log.Fatalln("could not start server:", err)
	}
}
