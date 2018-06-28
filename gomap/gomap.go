package gomap

import (
	"errors"

	"github.com/osmlab/gomap/db"
)

var (
	// ErrElementNotFound determines that element doesn't exist
	ErrElementNotFound = errors.New("element doesn't exist")
	// ErrElementDeleted determines that element is deleted
	ErrElementDeleted = errors.New("element is deleted")
)

// Gomap contains business logic of Openstreetmap server
type Gomap struct {
	db *db.OsmDB
}

// New returns new Gomap
func New(db *db.OsmDB) *Gomap {
	return &Gomap{db: db}
}
