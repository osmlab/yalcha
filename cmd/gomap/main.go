package main

import (
	"log"

	"github.com/osmlab/gomap/gomap"

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
	g := gomap.New(db)
	server := server.New(g)
	router := router.Load(config, server)
	err = router.Start(":" + config.Port)
	log.Fatalf("Server started with error: %v", err)
}
