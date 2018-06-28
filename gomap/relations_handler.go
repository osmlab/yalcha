package gomap

import "github.com/osmlab/gomap/osm"

// RelationsHandler is used to get data for /api/0.6/relations?relations=... request
func (g *Gomap) RelationsHandler(ids []int64, histIDs [][2]int64) (*osm.OSM, error) {
	current, err := g.db.ExtractRelations(ids)
	if err != nil {
		return nil, err
	}

	historic, err := g.db.ExtractHistoricalRelations(histIDs)
	if err != nil {
		return nil, err
	}

	relations := append(current, historic...)

	if len(relations) != len(ids)+len(histIDs) {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	resp.Relations = relations
	return resp, nil
}
