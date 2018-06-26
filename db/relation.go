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

// GetRelationFull is used to return relation with all internal members
func (o *OsmDB) GetRelationFull(id int64) (*osm.OSM, error) {
	query := fmt.Sprintf(`
	WITH relation AS (
		SELECT * FROM get_relation_by_id(%v)
	), relation_members AS (
		SELECT unnest(members) AS members FROM relation
	), relation_ids AS (
		SELECT array_agg((members::osm_member).ref) FROM relation_members where (members::osm_member).type = 'relation'
	), relations AS (
		SELECT * FROM get_relation_by_id(
			variadic (SELECT * FROM relation_ids)
		)
		union
		SELECT * FROM relation
	), way_ids AS (
		SELECT array_agg((members::osm_member).ref) FROM relation_members where (members::osm_member).type = 'way'
	), ways AS (
		SELECT * FROM get_way_by_id(
			variadic (SELECT * FROM way_ids)
		)
	), node_ids AS (
		SELECT array_agg(ref) FROM (
			SELECT (members::osm_member).ref FROM relation_members where (members::osm_member).type = 'node'
			union
			SELECT unnest(nodes) AS ref FROM ways
		) AS r
	), nodes AS (
		SELECT * FROM get_node_by_id(
			variadic (SELECT * FROM node_ids)
		)
	), relations_array AS (
		SELECT array_to_json(array_agg(r)) AS relations FROM relations r
	), ways_array AS (
		SELECT array_to_json(array_agg(w)) AS ways FROM ways w
	), nodes_array AS (
		SELECT array_to_json(array_agg(n)) AS nodes FROM nodes n
	)
	SELECT COALESCE(r.relations, '[]'), COALESCE(w.ways, '[]'), COALESCE(n.nodes, '[]')
	FROM relations_array r, ways_array w, nodes_array n
	`, id)

	osm := osm.New()
	err := o.db.QueryRow(query).Scan(
		&osm.Relations,
		&osm.Ways,
		&osm.Nodes,
	)
	return osm, err
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
