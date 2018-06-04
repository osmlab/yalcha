-- +goose Up
create or replace function get_ways_for_node(
    id bigint
)
    returns setof osm_way stable parallel safe language sql as $$
with current_way_nodes as (
  select way_id, node_id, sequence_id
  from current_way_nodes
  where node_id = $1
)
select
    n.id,
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
        from current_way_tags t
        where t.way_id = n.id
    ) as tags,
    (
        select array_agg(node_id
        order by sequence_id)
        from current_way_nodes
    ) as nodes
from current_ways n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.id in (select way_id from current_way_nodes) 
and n.visible = true
$$;

-- +goose Down
drop function if exists get_ways_for_node(bigint);
