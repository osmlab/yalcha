package db

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres adapter
	"github.com/osmlab/gomap/config"

	"github.com/jackc/pgx"
)

const (
	stmtVisibleNode                = "visible_node"
	stmtVisibleWay                 = "visible_way"
	stmtVisibleRelation            = "visible_relation"
	stmtSelectNodes                = "select_nodes"
	stmtSelectWays                 = "select_ways"
	stmtSelectRelations            = "select_relations"
	stmtSelectNodesHistory         = "select_nodes_history"
	stmtSelectWaysHistory          = "select_ways_history"
	stmtSelectRelationsHistory     = "select_relations_history"
	stmtSelectHistoricalNodes      = "select_historical_nodes"
	stmtSelectHistoricalWays       = "select_historical_ways"
	stmtSelectHistoricalRelations  = "select_historical_relations"
	stmtSelectNodesFromBbox        = "visible_node_in_bbox"
	stmtExtractNodes               = "extract_nodes"
	stmtExtractWays                = "extract_ways"
	stmtExtractRelations           = "extract_relations"
	stmtExtractHistoricNodes       = "extract_historic_nodes"
	stmtExtractHistoricWays        = "extract_historic_ways"
	stmtExtractHistoricRelations   = "extract_historic_relations"
	stmtNodesFromWays              = "nodes_from_ways"
	stmtNodesFromRelations         = "nodes_from_relations"
	stmtWaysFromNodes              = "ways_from_nodes"
	stmtWaysFromRelations          = "ways_from_relations"
	stmtRelationMembersOfRelations = "relation_members_of_relations"
	stmtRelationParentsOfNodes     = "relation_parents_of_nodes"
	stmtRelationParentsOfWays      = "relation_parents_of_ways"
	stmtRelationParentsOfRelations = "relation_parents_of_relations"
)

// OsmDB contains logic to deal with Openstreetmap database
type OsmDB struct {
	pool *pgx.ConnPool
	db   *sqlx.DB

	statements map[string]*pgx.PreparedStatement
}

// Init returns new database connection
func Init(config config.DB) (*OsmDB, error) {
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.Host,
			Port:     config.Port,
			User:     config.User,
			Password: config.Password,
			Database: config.DBName,
		},
	})
	if err != nil {
		return nil, err
	}

	sts, err := initStatements(pool)
	if err != nil {
		return nil, err
	}

	dbinfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	conn, err := sqlx.Open("postgres", dbinfo)
	return &OsmDB{
		pool:       pool,
		db:         conn,
		statements: sts,
	}, err
}

func initStatements(conn *pgx.ConnPool) (map[string]*pgx.PreparedStatement, error) {
	sts := make(map[string]*pgx.PreparedStatement)

	if _, err := conn.Prepare(
		stmtVisibleNode,
		strings.TrimSpace(`
			SELECT visible 
			FROM current_nodes 
			WHERE id = $1
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtVisibleWay,
		strings.TrimSpace(`
			SELECT visible 
			FROM current_ways 
			WHERE id = $1
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtVisibleRelation,
		strings.TrimSpace(`
			SELECT visible 
			FROM current_relations 
			WHERE id = $1
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectNodes,
		strings.TrimSpace(`
			SELECT id
			FROM current_nodes
			WHERE id = ANY($1)
		`),
	); err != nil {
		return sts, err
	}

	if _, err := conn.Prepare(
		stmtSelectWays,
		strings.TrimSpace(`
			SELECT id
			FROM current_ways
			WHERE id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectRelations,
		strings.TrimSpace(`
			SELECT id
			FROM current_relations
			WHERE id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectNodesHistory,
		strings.TrimSpace(`
			SELECT 
				node_id AS id, 
				version
			FROM nodes
			WHERE 
				node_id = ANY($1) AND
				(redaction_id IS NULL)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectWaysHistory,
		strings.TrimSpace(`
			SELECT 
				way_id AS id, 
				version
			FROM ways
			WHERE 
				way_id = ANY($1) AND
				(redaction_id IS NULL)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectRelationsHistory,
		strings.TrimSpace(`
			SELECT 
				relation_id AS id, 
				version
			FROM relations
			WHERE 
				relation_id = ANY($1) AND
				(redaction_id IS NULL)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectHistoricalNodes,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				n.node_id AS id, 
				n.version
			FROM nodes n
			INNER JOIN wanted w ON n.node_id = w.id AND
					   n.version = w.version
			WHERE (n.redaction_id IS NULL OR $3 = TRUE)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectHistoricalWays,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				w.way_id AS id, 
				w.version
			FROM ways w
			INNER JOIN wanted x ON w.way_id = x.id AND
					   w.version = x.version
			WHERE (w.redaction_id IS NULL OR $3 = TRUE)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectHistoricalRelations,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				r.relation_id AS id, 
				r.version
			FROM relations r
			INNER JOIN wanted x ON r.relation_id = x.id AND
					   r.version = x.version
			WHERE (r.redaction_id IS NULL OR $3 = TRUE)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtWaysFromNodes,
		strings.TrimSpace(`
			SELECT DISTINCT wn.way_id AS id
			FROM current_way_nodes wn
			WHERE wn.node_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtSelectNodesFromBbox,
		strings.TrimSpace(`
			SELECT id
			FROM current_nodes
			WHERE latitude BETWEEN $1 AND $2 AND 
				  longitude BETWEEN $3 AND $4 AND visible = true
			LIMIT $5
		`),
	); err != nil {
		return nil, err
	}

	_, err := conn.Prepare(
		stmtExtractNodes,
		strings.TrimSpace(`
			SELECT 
				n.id, 
				n.latitude / 1e7 :: float,
				n.longitude / 1e7 :: float,
				n.visible,
				n.timestamp,
				n.changeset_id,
				u.display_name,
				u.id,
				n.version,
				COALESCE(to_json((
					SELECT 
						array_agg(
							(k, v) :: osm_tag
							ORDER BY k
						)
					FROM current_node_tags t
					WHERE t.node_id = n.id
				)), '[]') as tags
			FROM current_nodes n
			JOIN changesets c ON c.id = n.changeset_id
    		LEFT JOIN users u ON (u.id = c.user_id and u.data_public)
			WHERE n.id = ANY($1)
			ORDER BY n.id
		`),
	)
	if err != nil {
		return nil, err
	}

	_, err = conn.Prepare(
		stmtExtractWays,
		strings.TrimSpace(`
			SELECT 
				w.id, 
				w.visible,
				w.timestamp,
				w.changeset_id,
				u.display_name,
    			u.id,
				w.version, 
				COALESCE(to_json(t.tags), '[]') as tags,
				COALESCE(to_json(wn.node_ids), '[]') as node_ids
			FROM current_ways w 
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(k, v) :: osm_tag
						ORDER BY k
					) as tags
				FROM current_way_tags 
				WHERE w.id=way_id
			) t ON true
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(node_id) as node_ids
				FROM (
					SELECT 
						node_id 
					FROM current_way_nodes 
					WHERE w.id=way_id
					ORDER BY sequence_id
				) x
			) wn ON true
			JOIN changesets c ON c.id = w.changeset_id
    		LEFT JOIN users u ON (u.id = c.user_id and u.data_public) 
			WHERE w.id = ANY($1)
			ORDER BY w.id
		`),
	)
	if err != nil {
		return nil, err
	}

	_, err = conn.Prepare(
		stmtExtractRelations,
		strings.TrimSpace(`
			SELECT 
				r.id, 
				r.visible, 
				r.timestamp,
				r.changeset_id,
				u.display_name,
				u.id,
				r.version, 
				COALESCE(to_json(t.tags), '[]') as tags,
				COALESCE(to_json(rm.members), '[]') AS members
			FROM current_relations r 
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(k, v) :: osm_tag
						ORDER BY k
					) as tags 
				FROM current_relation_tags 
				WHERE r.id=relation_id
			) t ON true
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(
							case member_type
							when 'Way'
								then 'way'
							when 'Relation'
								then 'relation'
							when 'Node'
								then 'node'
							end,
							member_id,
							member_role
						) :: osm_member
						ORDER BY sequence_id
					) as members
				FROM current_relation_members t
				WHERE t.relation_id = r.id
			) rm ON true
			JOIN changesets c ON c.id = r.changeset_id
    		LEFT JOIN users u ON (u.id = c.user_id and u.data_public)
			WHERE r.id = ANY($1)
			ORDER BY r.id
		`),
	)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtExtractHistoricNodes,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				n.node_id AS id, 
				n.latitude / 1e7 :: float,
				n.longitude / 1e7 :: float, 
				n.visible,
				n.timestamp,
				n.changeset_id, 
				u.display_name,
				u.id,
				n.version, 
				COALESCE(to_json((
					SELECT 
						array_agg(
							(k, v) :: osm_tag
							ORDER BY k
						)
					FROM node_tags t
					WHERE t.node_id = n.node_id AND
						  t.version = n.version
				)), '[]') as tags 
			FROM nodes n 
			INNER JOIN wanted x ON n.node_id = x.id AND 
					   n.version = x.version 
			JOIN changesets c ON c.id = n.changeset_id
			LEFT JOIN users u ON (u.id = c.user_id and u.data_public)
			ORDER BY n.node_id, n.version
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtExtractHistoricWays,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				w.way_id AS id, 
				w.visible,
				w.timestamp,
				w.changeset_id, 
				u.display_name,
    			u.id,
				w.version, 
				COALESCE(to_json(t.tags), '[]') as tags,
				COALESCE(to_json(wn.node_ids), '[]') as node_ids
			FROM ways w 
			INNER JOIN wanted x ON w.way_id = x.id AND 
								   w.version = x.version
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(k, v) :: osm_tag
						ORDER BY k
					) as tags
				FROM way_tags
				WHERE w.way_id=way_id AND
					  w.version=version
			) t ON true
			LEFT JOIN LATERAL (
				SELECT array_agg(node_id) as node_ids
				FROM (
					SELECT node_id 
					FROM way_nodes
					WHERE w.way_id=way_id AND 
						  w.version=version 
					ORDER BY sequence_id
				) x
			) wn ON true
			JOIN changesets c ON c.id = w.changeset_id
    		LEFT JOIN users u ON (u.id = c.user_id and u.data_public)
			ORDER BY w.way_id, w.version
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtExtractHistoricRelations,
		strings.TrimSpace(`
			WITH wanted(id, version) AS (
				SELECT * 
				FROM unnest(CAST($1 AS bigint[]), CAST($2 AS bigint[]))
			)
			SELECT 
				r.relation_id AS id, 
				r.visible,
				r.timestamp,
				r.changeset_id,
				u.display_name,
				u.id,
				r.version, 
				COALESCE(to_json(t.tags), '[]') as tags,
				COALESCE(to_json(rm.members), '[]') AS members
			FROM relations r 
			INNER JOIN wanted x ON r.relation_id = x.id AND
								   r.version = x.version
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(k, v) :: osm_tag
						ORDER BY k
					) as tags 
				FROM relation_tags 
				WHERE r.relation_id=relation_id AND
					  r.version=version
			) t ON true 
			LEFT JOIN LATERAL (
				SELECT 
					array_agg(
						(
							case member_type
							when 'Way'
								then 'way'
							when 'Relation'
								then 'relation'
							when 'Node'
								then 'node'
							end,
							member_id,
							member_role
						) :: osm_member
						ORDER BY sequence_id
					) as members
				FROM relation_members t
				WHERE t.relation_id = r.relation_id AND
					  t.version = r.version
			) rm ON true
			JOIN changesets c ON c.id = r.changeset_id
    		LEFT JOIN users u ON (u.id = c.user_id and u.data_public)
			ORDER BY r.relation_id, r.version
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtNodesFromWays,
		strings.TrimSpace(`
			SELECT DISTINCT wn.node_id AS id
			FROM current_way_nodes wn
			WHERE wn.way_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtNodesFromRelations,
		strings.TrimSpace(`
			SELECT DISTINCT rm.member_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Node' AND 
				  rm.relation_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtWaysFromNodes,
		strings.TrimSpace(`
			SELECT DISTINCT wn.way_id AS id
			FROM current_way_nodes wn
			WHERE wn.node_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtWaysFromRelations,
		strings.TrimSpace(`
			SELECT DISTINCT rm.member_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Way' AND 
				  rm.relation_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtRelationMembersOfRelations,
		strings.TrimSpace(`
			SELECT DISTINCT rm.member_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Relation' AND 
				  rm.relation_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtRelationParentsOfNodes,
		strings.TrimSpace(`
			SELECT DISTINCT rm.relation_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Node' AND 
				  rm.member_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtRelationParentsOfWays,
		strings.TrimSpace(`
			SELECT DISTINCT rm.relation_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Way' AND 
				  rm.member_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	if _, err := conn.Prepare(
		stmtRelationParentsOfRelations,
		strings.TrimSpace(`
			SELECT DISTINCT rm.relation_id AS id
			FROM current_relation_members rm
			WHERE rm.member_type = 'Relation' AND 
				  rm.member_id = ANY($1)
		`),
	); err != nil {
		return nil, err
	}

	return sts, nil
}
