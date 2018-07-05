package gomap

import (
	"github.com/osmlab/gomap/osm"
)

// RelationFullHandler is used to get data for /api/0.6/relation/.../full request
func (g *Gomap) RelationFullHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectRelations(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	for i := range ids {
		isVisible, err := g.db.IsRelationVisible(ids[i])
		if err != nil {
			return nil, err
		}
		if !isVisible {
			return nil, ErrElementDeleted
		}
	}

	nodesFromRelations, err := g.db.SelectNodesFromRelations(ids)
	if err != nil {
		return nil, err
	}
	waysFromRelations, err := g.db.SelectWaysFromRelations(ids)
	if err != nil {
		return nil, err
	}
	nodesFromWays, err := g.db.SelectNodesFromWays(waysFromRelations)
	if err != nil {
		return nil, err
	}
	relationsFromRelations, err := g.db.SelectRelationMembersFromRelations(ids)
	if err != nil {
		return nil, err
	}

	nodeIDs := append(nodesFromRelations, nodesFromWays...)
	wayIDs := waysFromRelations
	relationIDs := append(ids, relationsFromRelations...)

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
