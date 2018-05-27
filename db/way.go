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
	FROM get_way_by_id(%v)
	`, id)

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

// GetWayFull selects way with internal nodes from database by id
func (o *OsmDB) GetWayFull(id int64) (*osm.OSM, error) {
	query := fmt.Sprintf(`
	WITH way AS (
		SELECT * FROM get_way_by_id(%v)
	), node_ids AS (
		SELECT ARRAY_AGG(ref) from (
			SELECT UNNEST(nodes) AS ref FROM way
		) AS r
	), nodes AS (
		SELECT * FROM get_node_by_id(
			VARIADIC (SELECT * FROM node_ids)
		)
	), ways_array AS (
		SELECT array_to_json(array_agg(w)) AS ways FROM way w
	), nodes_array AS (
		SELECT array_to_json(array_agg(n)) AS nodes FROM nodes n
	)
	SELECT COALESCE(w.ways, '[]'), COALESCE(n.nodes, '[]')
	FROM ways_array w, nodes_array n
	`, id)

	osm := osm.New()
	err := o.db.QueryRow(query).Scan(
		&osm.Ways,
		&osm.Nodes,
	)
	return osm, err
}

// GetWays returns ways by ids
func (o *OsmDB) GetWays(ids []int64) (*osm.Ways, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	waysQuery := fmt.Sprintf(`
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
	`, arrayToString(ids))

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
