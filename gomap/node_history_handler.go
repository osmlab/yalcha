package gomap

import "github.com/osmlab/gomap/osm"

// NodeHistoryHandler is used to get data for /api/0.6/node/.../history request
func (g *Gomap) NodeHistoryHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectNodesHistory(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	nodes, err := g.db.ExtractHistoricalNodes(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Nodes = nodes
	return resp, nil
}
