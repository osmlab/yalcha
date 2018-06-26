package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/gomap/osm"
)

// GetRelationID selects relation id from database by id
func (o *OsmDB) GetRelationID(id int64) (int64, error) {
	var relationID int64
	err := o.pool.QueryRow(stmtSelectRelations, []int64{id}).Scan(&relationID)
	return relationID, err
}

// IsRelationVisible is used to check relation visibility
func (o *OsmDB) IsRelationVisible(id int64) (bool, error) {
	var result bool
	err := o.pool.QueryRow(stmtVisibleRelation, id).Scan(&result)
	return result, err
}

// GetRelations selects relations from database by id
func (o *OsmDB) GetRelations(ids []int64) (osm.Relations, error) {
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

// GetHistoricRelations selects historical relations from database by id and version
func (o *OsmDB) GetHistoricRelations(ids [][2]int64) (osm.Relations, error) {
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

// GetRelationsFromRelations selects relations id from database by id and relations ids
func (o *OsmDB) GetRelationsFromRelations(ids []int64) ([]int64, error) {
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

// GetRelationByVersion selects relation from database by id and version
func (o *OsmDB) GetRelationByVersion(id, version int64) (*osm.Relation, error) {
	relationQuery := fmt.Sprintf(`
	SELECT
		id,
		visible,
		version,
		"user",
		uid,
		changeset,
		timestamp,
		COALESCE(to_json(tags), '[]') AS tags,
		COALESCE(to_json(members), '[]') AS members
	FROM get_relation_by_id_and_version(array[[%v, %v]])
	`, id, version)

	var user sql.NullString
	var userID sql.NullInt64
	relation := &osm.Relation{}
	err := o.db.QueryRow(relationQuery).Scan(
		&relation.ID,
		&relation.Visible,
		&relation.Version,
		&user,
		&userID,
		&relation.ChangesetID,
		&relation.Timestamp,
		&relation.Tags,
		&relation.Members,
	)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		relation.UserID = &userID.Int64
		if user.Valid {
			relation.User = &user.String
		}
	}

	return relation, nil
}

// GetRelationHistory selects relation history from databASe by id
func (o *OsmDB) GetRelationHistory(id int64) (osm.Relations, error) {
	query := fmt.Sprintf(`
	SELECT
		id,
		visible,
		version,
		"user",
		uid,
		changeset,
		timestamp,
		COALESCE(to_json(tags), '[]') AS tags,
		COALESCE(to_json(members), '[]') AS members
	FROM get_relation_history_by_id(%v)`, id)

	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations osm.Relations
	for rows.Next() {
		var user sql.NullString
		var userID sql.NullInt64
		relation := &osm.Relation{}
		err := rows.Scan(
			&relation.ID,
			&relation.Visible,
			&relation.Version,
			&user,
			&userID,
			&relation.ChangesetID,
			&relation.Timestamp,
			&relation.Tags,
			&relation.Members,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			relation.UserID = &userID.Int64
			if user.Valid {
				relation.User = &user.String
			}
		}

		relations = append(relations, relation)
	}

	return relations, nil
}
