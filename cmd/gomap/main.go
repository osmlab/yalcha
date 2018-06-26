package main

import (
	"log"

	"github.com/osmlab/gomap/config"
	"github.com/osmlab/gomap/db"
	"github.com/osmlab/gomap/router"
	"github.com/osmlab/gomap/server"
)

func main() {
	config := &config.Config{
		Port: "8090",
		Database: config.DB{
			Host:     "localhost",
			Port:     5432,
			DBName:   "openstreetmap",
			User:     "heorhi",
			Password: "some_password",
		},
	}
	db, err := db.Init(config.Database)
	if err != nil {
		log.Fatalf("DB started with error: %v", err)
	}
	server := server.New(db)
	router := router.Load(config, server)
	err = router.Start(":" + config.Port)
	log.Fatalf("Server started with error: %v", err)
}
