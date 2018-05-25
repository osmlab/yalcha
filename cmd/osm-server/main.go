package main

import (
	"log"

	"github.com/osmlab/yalcha/config"
	"github.com/osmlab/yalcha/db"
	"github.com/osmlab/yalcha/router"
	"github.com/osmlab/yalcha/server"
)

func main() {
	config := &config.Config{
		Port: "8090",
	}
	db, err := db.Init()
	if err != nil {
		log.Fatalf("DB started with error: %v", err)
	}
	server := server.New(db)
	router := router.Load(config, server)
	err = router.Start("localhost:" + config.Port)
	log.Fatalf("Server started with error: %v", err)
}
