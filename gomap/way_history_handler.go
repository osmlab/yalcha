package gomap

import "github.com/osmlab/gomap/osm"

// WayHistoryHandler is used to get data for /api/0.6/way/.../history request
func (g *Gomap) WayHistoryHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectWaysHistory(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	ways, err := g.db.ExtractHistoricalWays(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Ways = ways
	return resp, nil
}
