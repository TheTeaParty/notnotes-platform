create table tags
(
    id         varchar                                               not null
        constraint tags_pk
            primary key,
    name       varchar                                               not null,
    created_at timestamp without time zone default current_timestamp not null
);

create unique index tags_id_uindex
    on tags (id);

create unique index tags_name_uindex
    on tags (name);

