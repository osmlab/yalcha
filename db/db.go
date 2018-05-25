package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres adapter
)

const (
	dbUser = "heorhi"
	dbName = "gis"
)

// OsmDB contains logic to deal with Openstreetmap database
type OsmDB struct {
	db *sqlx.DB
}

// Init returns new database connection
func Init() (*OsmDB, error) {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		dbUser, dbName)

	conn, err := sqlx.Open("postgres", dbinfo)
	return &OsmDB{db: conn}, err
}
