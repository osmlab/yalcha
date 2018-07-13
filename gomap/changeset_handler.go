package gomap

import (
	"github.com/osmlab/gomap/osm"
)

// ChangesetHandler is used to get data for /api/0.6/changeset/... request
func (g *Gomap) ChangesetHandler(id int64, includeDiscussion bool) (*osm.OSM, error) {
	ids, err := g.db.SelectChangesets(id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrElementNotFound
	}

	changesets, err := g.db.ExtractChangesets(ids, includeDiscussion)
	if err != nil {
		return nil, err
	}

	resp := osm.New()
	resp.Changesets = changesets
	return resp, nil
}
