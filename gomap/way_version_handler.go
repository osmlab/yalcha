package gomap

import "github.com/osmlab/gomap/osm"

// WayVersionHandler returns data for /api/0.6/way/:id/:version request
func (g *Gomap) WayVersionHandler(id, version int64) (*osm.OSM, error) {
	var histIDs [][2]int64
	histIDs = append(histIDs, [2]int64{id, version})

	ways, err := g.db.ExtractHistoricalWays(histIDs)
	if err != nil {
		return nil, err
	}
	if len(ways) == 0 {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	resp.Ways = ways
	return resp, nil
}
