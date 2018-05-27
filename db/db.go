package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres adapter
	"github.com/osmlab/yalcha/config"
)

// OsmDB contains logic to deal with Openstreetmap database
type OsmDB struct {
	db *sqlx.DB
}

// Init returns new database connection
func Init(config config.DB) (*OsmDB, error) {
	dbinfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	conn, err := sqlx.Open("postgres", dbinfo)
	return &OsmDB{db: conn}, err
}
