package db

import (
	"github.com/osmlab/gomap/osm"
)

// SelectWays selects ways id
func (o *OsmDB) SelectWays(ids ...int64) ([]int64, error) {
	var result []int64
	rows, err := o.pool.Query(stmtSelectWays, ids)
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

// IsWayVisible is used to check way visibility
func (o *OsmDB) IsWayVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleWay, id).Scan(&result)
	return result, err
}

// ExtractWays selects ways from database by id
func (o *OsmDB) ExtractWays(ids []int64) (osm.Ways, error) {
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

// ExtractHistoricalWays returns historical ways by id and version
func (o *OsmDB) ExtractHistoricalWays(ids [][2]int64) (osm.Ways, error) {
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

// SelectWaysFromRelations selects ways id from database by id and relations ids
func (o *OsmDB) SelectWaysFromRelations(ids []int64) ([]int64, error) {
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
