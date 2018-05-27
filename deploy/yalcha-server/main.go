package main

import (
	"log"
	"os"

	"github.com/osmlab/yalcha/config"
	"github.com/osmlab/yalcha/db"
	"github.com/osmlab/yalcha/router"
	"github.com/osmlab/yalcha/server"
)

var database *db.OsmDB

func init() {
	var err error
	config := config.DB{
		Host:     "gis.cesfozorknmw.us-west-2.rds.amazonaws.com",
		Port:     "5432",
		DBName:   "gis",
		User:     "hesidoryn",
		Password: "hesidoryn",
	}
	database, err = db.Init(config)
	if err != nil {
		log.Fatalf("DB started with error: %v", err)
	}
}

func main() {
	config := &config.Config{
		Port: os.Getenv("PORT"),
	}
	server := server.New(database)
	router := router.Load(config, server)
	err := router.Start(":" + config.Port)
	log.Fatalf("Server started with error: %v", err)
}
