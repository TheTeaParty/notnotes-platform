create table note_tags
(
    id         varchar                                               not null
        constraint note_tags_pk
            primary key,
    note_id    varchar                                               not null,
    tag_id     varchar                                               not null,
    created_at timestamp without time zone default current_timestamp not null
);

create unique index note_tags_id_uindex
    on note_tags (id);

