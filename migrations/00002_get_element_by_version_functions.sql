-- +goose Up
create or replace function get_node_by_id_and_version(ids_versions bigint[])
  returns setof osm_node stable parallel safe language sql as $$
select distinct
  n.node_id as id,
  n.visible,
  n.version,
  n.changeset_id,
  n.timestamp,
  u.display_name,
  u.id,
  n.latitude / 1e7 :: float,
  n.longitude / 1e7 :: float,
  --n.tags
  (
    select array_agg(
        (k, v) :: osm_tag
    order by k)
    from node_tags t
    where t.node_id = n.node_id
  ) as tags
from generate_subscripts($1, 1) i
  join   nodes n on n.node_id = $1[i][1] and n.version= $1[i][2]
  join changesets c on c.id = n.changeset_id
  left join users u on (u.id = c.user_id and u.data_public)
$$;

create or replace function get_way_by_id_and_version(ids_versions bigint[])
  returns setof osm_way stable parallel safe language sql as $$
select distinct
    n.way_id,
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
from generate_subscripts($1, 1) i
  join ways n on n.way_id = $1[i][1] and n.version= $1[i][2]
  join changesets c on c.id = n.changeset_id
  left join users u on (u.id = c.user_id and u.data_public)
$$;

create or replace function get_relation_by_id_and_version(ids_versions bigint[])
  returns setof osm_relation stable parallel safe language sql as $$
select distinct
    n.relation_id,
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
from generate_subscripts($1, 1) i
  join relations n on n.relation_id = $1[i][1] and n.version= $1[i][2]
  join changesets c on c.id = n.changeset_id
  left join users u on (u.id = c.user_id and u.data_public)
$$;

-- +goose Down

drop function if exists get_node_by_id_and_version (bigint[]);
drop function if exists get_way_by_id_and_version (bigint[]);
drop function if exists get_relation_by_id_and_version (bigint[]);