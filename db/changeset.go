package db

import (
	"github.com/osmlab/gomap/osm"
)

// SelectChangesets selects changesets id
func (o *OsmDB) SelectChangesets(ids ...int64) ([]int64, error) {
	var result []int64
	rows, err := o.pool.Query(stmtSelectChangesets, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

// ExtractChangesets extract changesets from database by id
func (o *OsmDB) ExtractChangesets(ids []int64, includeDiscussion bool) (osm.Changesets, error) {
	rows, err := o.pool.Query(stmtExtractChangesets, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	changesets := osm.Changesets{}
	for rows.Next() {
		changeset := &osm.Changeset{
			Discussion: &osm.ChangesetDiscussion{},
		}
		if err := rows.Scan(
			&changeset.ID,
			&changeset.UserID,
			&changeset.User,
			&changeset.CreatedAt,
			&changeset.ClosedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ChangesCount,
			&changeset.Tags,
			&changeset.CommentsCount,
			&changeset.Discussion.Comments,
		); err != nil {
			return nil, err
		}

		if !includeDiscussion {
			changeset.Discussion = nil
		}

		changesets = append(changesets, changeset)
	}

	return changesets, nil
}
