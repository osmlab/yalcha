package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/gomap/osm"
)

// GetWayID selects way id from database by id
func (o *OsmDB) GetWayID(id int64) (int64, error) {
	var wayID int64
	err := o.pool.QueryRow(stmtSelectWays, []int64{id}).Scan(&wayID)
	return wayID, err
}

// IsWayVisible is used to check way visibility
func (o *OsmDB) IsWayVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleWay, id).Scan(&result)
	return result, err
}

// GetWays selects ways from database by id
func (o *OsmDB) GetWays(ids []int64) (osm.Ways, error) {
	rows, err := o.pool.Query(stmtExtractWays, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ways := osm.Ways{}
	for rows.Next() {
		way := &osm.Way{}
		if err := rows.Scan(
			&way.ID,
			&way.Visible,
			&way.Timestamp,
			&way.ChangesetID,
			&way.User,
			&way.UserID,
			&way.Version,
			&way.Tags,
			&way.Nodes,
		); err != nil {
			return nil, err
		}
		ways = append(ways, way)
	}

	return ways, nil
}

// GetHistoricWays selects historical ways from database by id and version
func (o *OsmDB) GetHistoricWays(ids [][2]int64) (osm.Ways, error) {
	wayIDs, vers := []int64{}, []int64{}
	for i := range ids {
		wayIDs = append(wayIDs, ids[i][0])
		vers = append(vers, ids[i][1])
	}

	rows, err := o.pool.Query(stmtExtractHistoricWays, wayIDs, vers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ways := osm.Ways{}
	for rows.Next() {
		way := &osm.Way{}
		if err := rows.Scan(
			&way.ID,
			&way.Visible,
			&way.Timestamp,
			&way.ChangesetID,
			&way.User,
			&way.UserID,
			&way.Version,
			&way.Tags,
			&way.Nodes,
		); err != nil {
			return nil, err
		}
		ways = append(ways, way)
	}

	return ways, nil
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

// GetWaysFromRelations selects ways id from database by id and relations ids
func (o *OsmDB) GetWaysFromRelations(ids []int64) ([]int64, error) {
	rows, err := o.pool.Query(stmtWaysFromRelations, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wayIDs := []int64{}
	for rows.Next() {
		var wayID int64
		if err := rows.Scan(&wayID); err != nil {
			return nil, err
		}
		wayIDs = append(wayIDs, wayID)
	}

	return wayIDs, nil
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
