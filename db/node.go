package db

import (
	"github.com/osmlab/gomap/osm"
)

// SelectNodes selects nodes id
func (o *OsmDB) SelectNodes(ids ...int64) ([]int64, error) {
	var result []int64
	rows, err := o.pool.Query(stmtSelectNodes, ids)
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

// IsNodeVisible is used to check node visibility
func (o *OsmDB) IsNodeVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleNode, id).Scan(&result)
	return result, err
}

// ExtractNodes extract nodes from database by id
func (o *OsmDB) ExtractNodes(ids []int64) (osm.Nodes, error) {
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

// ExtractHistoricalNodes returns historical nodes by id and version
func (o *OsmDB) ExtractHistoricalNodes(ids [][2]int64) (osm.Nodes, error) {
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

// SelectNodesFromWays returns nodes id by ways ids
func (o *OsmDB) SelectNodesFromWays(ids []int64) ([]int64, error) {
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

// SelectNodesFromRelations selects nodes id from database by id and relations ids
func (o *OsmDB) SelectNodesFromRelations(ids []int64) ([]int64, error) {
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
