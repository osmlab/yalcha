package gomap

import "github.com/osmlab/gomap/osm"

// NodeWaysHandler is used to get data for /api/0.6/node/.../ways request
func (g *Gomap) NodeWaysHandler(id int64) (*osm.OSM, error) {
	ids, err := g.db.SelectNodes(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	wayIDs, err := g.db.SelectWaysFromNodes(ids...)
	if err != nil {
		return nil, err
	}
	if len(wayIDs) == 0 {
		return nil, ErrElementNotFound
	}

	for i := range wayIDs {
		isVisible, err := g.db.IsWayVisible(wayIDs[i])
		if err != nil {
			return nil, err
		}
		if !isVisible {
			return nil, ErrElementDeleted
		}
	}

	ways, err := g.db.ExtractWays(ids)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Ways = ways
	return resp, nil
}
