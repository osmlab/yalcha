package gomap

import (
	"github.com/osmlab/gomap/osm"
)

// RelationHandler is used to get date for /api/0.6/relation/... request
func (g *Gomap) RelationHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectRelations(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	for i := range ids {
		isVisible, err := g.db.IsRelationVisible(ids[i])
		if err != nil {
			return nil, err
		}
		if !isVisible {
			return nil, ErrElementDeleted
		}
	}

	relations, err := g.db.ExtractRelations(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Relations = relations
	return resp, nil
}
