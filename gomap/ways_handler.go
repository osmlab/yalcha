package gomap

import "github.com/osmlab/gomap/osm"

// WaysHandler is used to get data for /api/0.6/ways?ways=... request
func (g *Gomap) WaysHandler(ids []int64, histIDs [][2]int64) (*osm.OSM, error) {
	current, err := g.db.ExtractWays(ids)
	if err != nil {
		return nil, err
	}

	historic, err := g.db.ExtractHistoricalWays(histIDs)
	if err != nil {
		return nil, err
	}

	ways := append(current, historic...)

	if len(ways) != len(ids)+len(histIDs) {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	resp.Ways = ways
	return resp, nil
}
