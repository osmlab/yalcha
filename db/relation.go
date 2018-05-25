package db

import (
	"database/sql"
	"fmt"

	"github.com/osmlab/yalcha/osm"
)

// GetRelation selects relation from database by id
func (o *OsmDB) GetRelation(id int64) (*osm.Relation, error) {
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
FROM get_relation_by_id(%v)`, id)

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

// GetRelations returns relations by ids
func (o *OsmDB) GetRelations(ids []int64) (*osm.Relations, error) {
	relationsQuery := ""
	for i := range ids {
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
FROM get_relation_by_id(%v)
`, ids[i])
		relationsQuery += relationQuery
		if i != len(ids)-1 {
			relationsQuery += "UNION ALL"
		}
	}

	rows, err := o.db.Query(relationsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relations := &osm.Relations{}
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

		*relations = append(*relations, relation)
	}

	return relations, nil
}
