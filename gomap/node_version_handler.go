package gomap

import "github.com/osmlab/gomap/osm"

// NodeVersionHandler returns data for /api/0.6/node/:id/:version request
func (g *Gomap) NodeVersionHandler(id, version int64) (*osm.OSM, error) {
	var histIDs [][2]int64
	histIDs = append(histIDs, [2]int64{id, version})

	nodes, err := g.db.ExtractHistoricalNodes(histIDs)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	resp.Nodes = nodes
	return resp, nil
}
