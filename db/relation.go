package db

import (
	"github.com/osmlab/gomap/osm"
)

// SelectRelations selects relations ids
func (o *OsmDB) SelectRelations(ids ...int64) ([]int64, error) {
	var result []int64
	rows, err := o.pool.Query(stmtSelectRelations, ids)
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

// SelectRelationsHistory selects relations ids
func (o *OsmDB) SelectRelationsHistory(ids ...int64) ([][2]int64, error) {
	var result [][2]int64
	rows, err := o.pool.Query(stmtSelectRelationsHistory, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id [2]int64
		if err := rows.Scan(&id[0], &id[1]); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

// IsRelationVisible is used to check relation visibility
func (o *OsmDB) IsRelationVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleRelation, id).Scan(&result)
	return result, err
}

// ExtractRelations extract relations by id
func (o *OsmDB) ExtractRelations(ids []int64) (osm.Relations, error) {
	rows, err := o.pool.Query(stmtExtractRelations, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relations := osm.Relations{}
	for rows.Next() {
		relation := &osm.Relation{}
		if err := rows.Scan(
			&relation.ID,
			&relation.Visible,
			&relation.Timestamp,
			&relation.ChangesetID,
			&relation.User,
			&relation.UserID,
			&relation.Version,
			&relation.Tags,
			&relation.Members,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, err
}

// ExtractHistoricalRelations extarct historical relations from database by id and version
func (o *OsmDB) ExtractHistoricalRelations(ids [][2]int64) (osm.Relations, error) {
	relIDs, vers := []int64{}, []int64{}
	for i := range ids {
		relIDs = append(relIDs, ids[i][0])
		vers = append(vers, ids[i][1])
	}

	rows, err := o.pool.Query(stmtExtractHistoricRelations, relIDs, vers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relations := osm.Relations{}
	for rows.Next() {
		relation := &osm.Relation{}
		if err := rows.Scan(
			&relation.ID,
			&relation.Visible,
			&relation.Timestamp,
			&relation.ChangesetID,
			&relation.User,
			&relation.UserID,
			&relation.Version,
			&relation.Tags,
			&relation.Members,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, err
}

// SelectRelationsFromRelations selects relations id from database by id and relations ids
func (o *OsmDB) SelectRelationsFromRelations(ids []int64) ([]int64, error) {
	rows, err := o.pool.Query(stmtRelationsFromRelations, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relIDs := []int64{}
	for rows.Next() {
		var relID int64
		if err := rows.Scan(&relID); err != nil {
			return nil, err
		}
		relIDs = append(relIDs, relID)
	}

	return relIDs, nil
}
