package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/yalcha/osm"
)

// GetWay selects way from database by id
func (o *OsmDB) GetWay(id int64) (*osm.Way, error) {
	wayQuery := fmt.Sprintf(`
SELECT
    id, 
    visible, 
    version,
    "user",
    uid,
    changeset, 
    timestamp,
    COALESCE(to_json(nodes), '[]') AS nodes,
    COALESCE(to_json(tags), '[]') AS tags
FROM get_way_by_id(%v)`, id)

	var user sql.NullString
	var userID sql.NullInt64
	way := &osm.Way{}
	err := o.db.QueryRow(wayQuery).Scan(
		&way.ID,
		&way.Visible,
		&way.Version,
		&user,
		&userID,
		&way.ChangesetID,
		&way.Timestamp,
		&way.Nodes,
		&way.Tags,
	)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		way.UserID = &userID.Int64
		if user.Valid {
			way.User = &user.String
		}
	}

	return way, nil
}

// GetWays returns ways by ids
func (o *OsmDB) GetWays(ids []int64) (*osm.Ways, error) {
	waysQuery := ""
	for i := range ids {
		wayQuery := fmt.Sprintf(`
SELECT
	id, 
	visible, 
	version,
	"user",
	uid,
	changeset, 
	timestamp,
	COALESCE(to_json(nodes), '[]') AS nodes,
	COALESCE(to_json(tags), '[]') AS tags
FROM get_way_by_id(%v)
`, ids[i])
		waysQuery += wayQuery
		if i != len(ids)-1 {
			waysQuery += "UNION ALL"
		}
	}

	rows, err := o.db.Query(waysQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ways := &osm.Ways{}
	for rows.Next() {
		var user sql.NullString
		var userID sql.NullInt64
		way := &osm.Way{}
		err := rows.Scan(
			&way.ID,
			&way.Visible,
			&way.Version,
			&user,
			&userID,
			&way.ChangesetID,
			&way.Timestamp,
			&way.Nodes,
			&way.Tags,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			way.UserID = &userID.Int64
			if user.Valid {
				way.User = &user.String
			}
		}

		*ways = append(*ways, way)
	}

	return ways, nil
}
