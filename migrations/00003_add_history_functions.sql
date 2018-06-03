-- +goose Up
create or replace function get_node_history_by_id(
    id bigint
)
    returns setof osm_node stable parallel safe language sql as $$
select
    n.node_id as id,
    n.visible,
    n.version,
    n.changeset_id,
    n.timestamp,
    u.display_name,
    u.id,
    n.latitude / 1e7 :: float,
    n.longitude / 1e7 :: float,
    (
        select array_agg(
            (k, v) :: osm_tag
        order by k)
        from node_tags t
        where t.node_id = n.node_id
    ) as tags
from nodes n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.node_id = $1
$$;

create or replace function get_way_history_by_id(
    id bigint
)
    returns setof osm_way stable parallel safe language sql as $$
select
    n.way_id as id,
    n.visible,
    n.version,
    n.changeset_id,
    n.timestamp,
    u.display_name,
    u.id,
    (
        select array_agg(
            (k, v) :: osm_tag
        order by k)
        from way_tags t
        where t.way_id = n.way_id
    ) as tags,
    (
        select array_agg(node_id
        order by sequence_id)
        from way_nodes t
        where t.way_id = n.way_id
    ) as nodes
from ways n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.way_id = $1
$$;

create or replace function get_relation_history_by_id(
    id bigint
)
    returns setof osm_relation stable parallel safe language sql as $$
select
    n.relation_id as id,
    n.visible,
    n.version,
    n.changeset_id,
    n.timestamp,
    u.display_name,
    u.id,
    (
        select array_agg(
            (k, v) :: osm_tag
        order by k)
        from relation_tags t
        where t.relation_id = n.relation_id
    ),
    (
        select array_agg(
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
        order by sequence_id)
        from relation_members t
        where t.relation_id = n.relation_id
    )
from relations n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.relation_id = $1
$$;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
