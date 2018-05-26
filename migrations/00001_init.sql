-- +goose Up

drop type if exists osm_node cascade;
drop type if exists osm_tag cascade;
drop type if exists osm_way cascade;
drop type if exists osm_relation cascade;
drop type if exists osm_member;

create type osm_tag as (
    k text,
    v text
);
create type osm_member as (
    type text,
    ref  bigint,
    role text
);
create type osm_node as (
    id          bigint,
    visible     boolean,
    "version"   bigint,
    changeset   bigint,
    "timestamp" timestamp,
    "user"      text,
    uid         bigint,
    lat         float,
    lon         float,
    tags        osm_tag []
);
create type osm_way as (
    id          bigint,
    visible     boolean,
    "version"   bigint,
    changeset   bigint,
    "timestamp" timestamp,
    "user"      text,
    uid         bigint,
    tags        osm_tag [],
    nodes       bigint []
);
create type osm_relation as (
    id          bigint,
    visible     boolean,
    "version"   bigint,
    changeset   bigint,
    "timestamp" timestamp,
    "user"      text,
    uid         bigint,
    tags        osm_tag [],
    members     osm_member []
);

create or replace function get_node_by_id(
    variadic p_ids bigint[]
)
    returns setof osm_node stable parallel safe language sql as $$
select
    n.id as id,
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
        from current_node_tags t
        where t.node_id = n.id
    ) as tags
from current_nodes n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.id = ANY(p_ids)
$$;

create or replace function get_way_by_id(
    variadic p_ids bigint[]
)
    returns setof osm_way stable parallel safe language sql as $$
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
        from current_way_nodes t
        where t.way_id = n.id
    ) as nodes
from current_ways n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.id = ANY(p_ids)
$$;
create or replace function get_relation_by_id(
    variadic p_ids bigint[]
)
    returns setof osm_relation stable parallel safe language sql as $$
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
        from current_relation_tags t
        where t.relation_id = n.id
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
        from current_relation_members t
        where t.relation_id = n.id
    )
from current_relations n
    join changesets c on c.id = n.changeset_id
    left join users u on (u.id = c.user_id and u.data_public)
where n.id = ANY(p_ids)
$$;

-- +goose Down

drop type osm_relation cascade;
drop type osm_way cascade;
drop type osm_node cascade;
drop type osm_tag;
drop type osm_member;

drop function if exists get_node_by_id( bigint );
drop function if exists get_way_by_id( bigint );
drop function if exists get_relation_by_id( bigint );