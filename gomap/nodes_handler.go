package gomap

import "github.com/osmlab/gomap/osm"

// NodesHandler is used to get data or /api/0.6/nodes?nodes=... request
func (g *Gomap) NodesHandler(ids []int64, histIDs [][2]int64) (*osm.OSM, error) {
	current, err := g.db.ExtractNodes(ids)
	if err != nil {
		return nil, err
	}

	historic, err := g.db.ExtractHistoricalNodes(histIDs)
	if err != nil {
		return nil, err
	}

	nodes := append(current, historic...)

	if len(nodes) != len(ids)+len(histIDs) {
		return nil, ErrElementNotFound
	}

	resp := osm.New()
	for i := range nodes {
		if !nodes[i].Visible {
			nodes[i].Lat = nil
			nodes[i].Lon = nil
		}
		resp.Nodes = append(resp.Nodes, nodes[i])
	}
	return resp, nil
}
