package gomap

import (
	"github.com/osmlab/gomap/osm"
)

// NodeHandler is used to get data for /api/0.6/node/... request
func (g *Gomap) NodeHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectNodes(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	for i := range ids {
		isVisible, err := g.db.IsNodeVisible(ids[i])
		if err != nil {
			return nil, err
		}
		if !isVisible {
			return nil, ErrElementDeleted
		}
	}

	nodes, err := g.db.ExtractNodes(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Nodes = nodes
	return resp, nil
}
