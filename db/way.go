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

// GetWayByVersion selects way from database by id and version
func (o *OsmDB) GetWayByVersion(id, version int64) (*osm.Way, error) {
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
	FROM get_way_by_id_and_version(array[[%v, %v]])
	`, id, version)

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

// GetWayHistory selects way history from database by id
func (o *OsmDB) GetWayHistory(id int64) (osm.Ways, error) {
	query := fmt.Sprintf(`
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
	FROM get_way_history_by_id(%v)
	`, id)

	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ways osm.Ways
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

		ways = append(ways, way)
	}

	return ways, nil
}

// GetWays returns ways by ids
func (o *OsmDB) GetWays(ids []int64, idvs [][2]int64) (osm.Ways, error) {
	if len(ids) == 0 &&
		len(idvs) == 0 {
		return nil, nil
	}

	query := ""
	if len(ids) > 0 {
		query = fmt.Sprintf(`
		SELECT
			id,
			visible,
			version,
			"user",
			uid,
			changeset,
			timestamp,
			COALESCE(to_jsonb(nodes), '[]') AS nodes,
			COALESCE(to_jsonb(tags), '[]') AS tags
		FROM get_way_by_id(%v)
		`, arrayToString(ids))
	}
	if len(idvs) > 0 {
		if len(query) > 0 {
			query += "UNION"
		}
		query += fmt.Sprintf(`
		SELECT
			id,
			visible,
			version,
			"user",
			uid,
			changeset,
			timestamp,
			COALESCE(to_jsonb(nodes), '[]') AS nodes,
			COALESCE(to_jsonb(tags), '[]') AS tags
		FROM get_way_by_id_and_version(array[[%v]])
		`, versionArrayToString(idvs))
	}

	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ways osm.Ways
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

		ways = append(ways, way)
	}

	return ways, nil
}
