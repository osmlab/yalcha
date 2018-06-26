package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/gomap/osm"
)

// GetNodeID selects node id from database by id
func (o *OsmDB) GetNodeID(id int64) (int64, error) {
	var nodeID int64
	err := o.pool.QueryRow(stmtSelectNodes, []int64{id}).Scan(&nodeID)
	return nodeID, err
}

// IsNodeVisible is used to check node visibility
func (o *OsmDB) IsNodeVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleNode, id).Scan(&result)
	return result, err
}

// GetNodes selects nodes from database by id
func (o *OsmDB) GetNodes(ids []int64) (osm.Nodes, error) {
	rows, err := o.pool.Query(stmtExtractNodes, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := osm.Nodes{}
	for rows.Next() {
		node := &osm.Node{}
		if err := rows.Scan(
			&node.ID,
			&node.Lat,
			&node.Lon,
			&node.Visible,
			&node.Timestamp,
			&node.ChangesetID,
			&node.User,
			&node.UserID,
			&node.Version,
			&node.Tags,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// GetHistoricalNodes selects historical nodes from database by id and version
func (o *OsmDB) GetHistoricalNodes(ids [][2]int64) (osm.Nodes, error) {
	nodeIDs, vers := []int64{}, []int64{}
	for i := range ids {
		nodeIDs = append(nodeIDs, ids[i][0])
		vers = append(vers, ids[i][1])
	}

	rows, err := o.pool.Query(stmtExtractHistoricNodes, nodeIDs, vers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := osm.Nodes{}
	for rows.Next() {
		node := &osm.Node{}
		if err := rows.Scan(
			&node.ID,
			&node.Lat,
			&node.Lon,
			&node.Visible,
			&node.Timestamp,
			&node.ChangesetID,
			&node.User,
			&node.UserID,
			&node.Version,
			&node.Tags,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// GetNodesFromWays selects nodes id from database by id and ways ids
func (o *OsmDB) GetNodesFromWays(ids []int64) ([]int64, error) {
	rows, err := o.pool.Query(stmtNodesFromWays, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodeIDs := []int64{}
	for rows.Next() {
		var nodeID int64
		if err := rows.Scan(&nodeID); err != nil {
			return nil, err
		}
		nodeIDs = append(nodeIDs, nodeID)
	}

	return nodeIDs, nil
}

// GetNodesFromRelations selects nodes id from database by id and relations ids
func (o *OsmDB) GetNodesFromRelations(ids []int64) ([]int64, error) {
	rows, err := o.pool.Query(stmtNodesFromRelations, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodeIDs := []int64{}
	for rows.Next() {
		var nodeID int64
		if err := rows.Scan(&nodeID); err != nil {
			return nil, err
		}
		nodeIDs = append(nodeIDs, nodeID)
	}

	return nodeIDs, nil
}

// GetNodeByVersion selects node from database by id and version
func (o *OsmDB) GetNodeByVersion(id, version int64) (*osm.Node, error) {
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
FROM get_node_by_id_and_version(array[[%v, %v]])`, id, version)

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

// GetNodeHistory selects node history from database by id
func (o *OsmDB) GetNodeHistory(id int64) (osm.Nodes, error) {
	query := fmt.Sprintf(`
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
	FROM get_node_history_by_id(%v)
	`, id)

	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes osm.Nodes
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

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// GetWaysForNode returns all the (not deleted) ways in which the given node is used
func (o *OsmDB) GetWaysForNode(id int64) (*osm.Ways, error) {
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
	FROM get_ways_for_node(%v)
	`, id)

	rows, err := o.db.Query(query)
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
