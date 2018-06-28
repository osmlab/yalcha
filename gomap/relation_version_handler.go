package gomap

import "github.com/osmlab/gomap/osm"

// RelationVersionHandler returns data for /api/0.6/relation/:id/:version request
func (g *Gomap) RelationVersionHandler(id, version int64) (*osm.OSM, error) {
	var histIDs [][2]int64
	histIDs = append(histIDs, [2]int64{id, version})

	relations, err := g.db.ExtractHistoricalRelations(histIDs)
	if err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	resp.Relations = relations
	return resp, nil
}
