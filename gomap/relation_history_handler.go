package gomap

import "github.com/osmlab/gomap/osm"

// RelationHistoryHandler is used to get data for /api/0.6/relation/.../history request
func (g *Gomap) RelationHistoryHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectRelationsHistory(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	relations, err := g.db.ExtractHistoricalRelations(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Relations = relations
	return resp, nil
}
