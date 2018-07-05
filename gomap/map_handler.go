package gomap

import (
	"github.com/osmlab/gomap/osm"
)

const maxNodes = 50000

// MapHandler is used to get data for /api/0.6/map?... request
func (g *Gomap) MapHandler(bbox []float64) (*osm.OSM, error) {
	nodesFromBbox, err := g.db.SelectNodesFromBbox(bbox)
	if err != nil {
		return nil, err
	}
	if len(nodesFromBbox) == 0 {
		return nil, ErrElementNotFound
	}
	if len(nodesFromBbox) > maxNodes {
		return nil, ErrElementNotFound
	}

	waysFromNodes, err := g.db.SelectWaysFromNodes(nodesFromBbox...)
	if err != nil {
		return nil, err
	}
	nodesFromWays, err := g.db.SelectNodesFromWays(waysFromNodes)
	if err != nil {
		return nil, err
	}
	relationsFromWays, err := g.db.SelectRelationsFromWays(waysFromNodes)
	if err != nil {
		return nil, err
	}
	relationsFromNodes, err := g.db.SelectRelationsFromNodes(append(nodesFromBbox, nodesFromWays...))
	if err != nil {
		return nil, err
	}
	relationsFromRelations, err := g.db.SelectRelationsFromRelations(append(relationsFromWays, relationsFromNodes...))
	if err != nil {
		return nil, err
	}

	nodeIDs := append(nodesFromBbox, nodesFromWays...)
	wayIDs := append(waysFromNodes)
	relationIDs := append(relationsFromWays, append(relationsFromRelations, relationsFromNodes...)...)

	nodes, err := g.db.ExtractNodes(nodeIDs)
	if err != nil {
		return nil, err
	}
	ways, err := g.db.ExtractWays(wayIDs)
	if err != nil {
		return nil, err
	}
	relations, err := g.db.ExtractRelations(relationIDs)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Nodes = nodes
	resp.Ways = ways
	resp.Relations = relations
	return resp, nil
}
