package gomap

import (
	"github.com/osmlab/gomap/osm"
)

// WayFullHandler is used to get data for /api/0.6/way/.../full request
func (g *Gomap) WayFullHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectWays(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	for i := range ids {
		isVisible, err := g.db.IsWayVisible(ids[i])
		if err != nil {
			return nil, err
		}
		if !isVisible {
			return nil, ErrElementDeleted
		}
	}

	nodeIDs, err := g.db.SelectNodesFromWays(ids)
	if err != nil {
		return nil, err
	}

	nodes, err := g.db.ExtractNodes(nodeIDs)
	if err != nil {
		return nil, err
	}
	ways, err := g.db.ExtractWays(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Nodes = nodes
	resp.Ways = ways
	return resp, nil
}
