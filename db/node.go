package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/yalcha/osm"
)

// GetNode selects node from database by id
func (o *OsmDB) GetNode(id int64) (*osm.Node, error) {
	nodeQuery := fmt.Sprintf(`
SELECT
	id, 
	visible, 
	version, 
	lat, 
	lon, 
	changeset, 
	"user",
	uid, 
	timestamp,
	COALESCE(to_json(tags), '[]') AS tags
FROM get_node_by_id(%v)`, id)

	var user sql.NullString
	var userID sql.NullInt64
	node := &osm.Node{}
	err := o.db.QueryRow(nodeQuery).Scan(
		&node.ID,
		&node.Visible,
		&node.Version,
		&node.Lat,
		&node.Lon,
		&node.ChangesetID,
		&user,
		&userID,
		&node.Timestamp,
		&node.Tags,
	)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		node.UserID = &userID.Int64
		if user.Valid {
			node.User = &user.String
		}
	}

	return node, nil
}

// GetNodes selects nodes from database by ids
func (o *OsmDB) GetNodes(ids []int64) (*osm.Nodes, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	nodesQuery := fmt.Sprintf(`
SELECT
	id, 
	visible, 
	version, 
	lat, 
	lon, 
	changeset, 
	"user",
	uid, 
	timestamp,
	COALESCE(to_json(tags), '[]') AS tags
FROM get_node_by_id(%v)
`, arrayToString(ids))

	rows, err := o.db.Query(nodesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := &osm.Nodes{}
	for rows.Next() {
		var user sql.NullString
		var userID sql.NullInt64
		node := &osm.Node{}
		err := rows.Scan(
			&node.ID,
			&node.Visible,
			&node.Version,
			&node.Lat,
			&node.Lon,
			&node.ChangesetID,
			&user,
			&userID,
			&node.Timestamp,
			&node.Tags,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			node.UserID = &userID.Int64
			if user.Valid {
				node.User = &user.String
			}
		}

		*nodes = append(*nodes, node)
	}

	return nodes, nil
}
